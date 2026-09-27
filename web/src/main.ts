const input = document.getElementById("inventoryInput") as HTMLTextAreaElement;
const button = document.getElementById("submitButton") as HTMLButtonElement;
const result = document.getElementById("result") as HTMLDivElement;
const test = document.getElementById("test") as HTMLPreElement;
const check = document.getElementById("jsonCheck") as HTMLInputElement;
const statusText = document.getElementById("status") as HTMLSpanElement;
const summary = document.getElementById("summary") as HTMLDivElement;
const itemResults = document.getElementById("itemResults") as HTMLDivElement;
const itemList = document.getElementById("itemList") as HTMLDivElement;
const authForm = document.getElementById("authForm") as HTMLFormElement;
const usernameInput = document.getElementById("usernameInput") as HTMLInputElement;
const passwordInput = document.getElementById("passwordInput") as HTMLInputElement;
const loginButton = document.getElementById("loginButton") as HTMLButtonElement;
const registerButton = document.getElementById("registerButton") as HTMLButtonElement;
const authError = document.getElementById("authError") as HTMLParagraphElement;
const authStatus = document.getElementById("authStatus") as HTMLDivElement;
const currentUser = document.getElementById("currentUser") as HTMLElement;
const logoutButton = document.getElementById("logoutButton") as HTMLButtonElement;

type MissingSet = {
  setName: string;
  missingParts: string[];
  missingCount: number;
};

type InventoryResult = {
  name: string;
  vaulted: boolean;
  ducats: number;
};

function renderMissingSets(sets: MissingSet[]) {
  summary.hidden = false;
  const missingPartCount = sets.reduce((count, set) => count + set.missingCount, 0);
  summary.innerHTML = `
    <div class="summary-item"><span class="summary-number">${sets.length}</span><span class="summary-label">incomplete sets</span></div>
    <div class="summary-item"><span class="summary-number">${missingPartCount}</span><span class="summary-label">missing parts</span></div>
  `;

  if (sets.length === 0) {
    result.className = "empty";
    result.textContent = "No incomplete Prime sets found in this inventory.";
    return;
  }

  result.className = "missing-grid";
  result.innerHTML = sets.map(set => `
    <article class="set-card">
      <h3>${escapeHtml(set.setName)}</h3>
      <ul class="part-list">
        ${set.missingParts.map(part => `<li>${escapeHtml(part)}</li>`).join("")}
      </ul>
    </article>
  `).join("");
}

function renderInventoryItems(items: InventoryResult[]) {
  if (items.length === 0) {
    itemResults.hidden = true;
    return;
  }

  itemResults.hidden = false;
  itemList.innerHTML = items.map(item => `
    <div class="item-row">
      <span>${escapeHtml(item.name)}</span>
      <span class="${item.vaulted ? "vaulted" : "available"}">${item.vaulted ? `Vaulted · ${item.ducats} ducats` : `Unvaulted · ${item.ducats} ducats`}</span>
    </div>
  `).join("");
}

function escapeHtml(value: string) {
  const element = document.createElement("div");
  element.textContent = value;
  return element.innerHTML;
}

button.addEventListener("click", async () => {
  const inventory = input.value.trim()

  let json
  if (!check.checked) {
    const items = inventory.split("\n")
      .map(line => line.trim())
      .filter(line => line.length > 0)
      .map(line => {
        const [name, quantity] = line.split(",").map(item => item.trim())
        const parsedQuantity = quantity ? parseInt(quantity, 10) : 1
        return { name, quantity: parsedQuantity }
      })
    json = JSON.stringify({ items })
  } else {
    json = inventory
  }

  try {
    test.textContent = (check.checked ? "json\n" : "nojson\n") + json;

    button.disabled = true;
    button.textContent = "Checking...";
    statusText.textContent = "Checking inventory...";
    summary.hidden = true;
    itemResults.hidden = true;
    result.className = "empty";
    result.textContent = "Loading results...";

    const response = await fetch("/inventory", {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: json
    });

    const text = await response.text();

    const lines = text
      .split("\n")
      .map(line => line.trim())
      .filter(line => line.length > 0);

    const data = lines.map(line => JSON.parse(line));
    const missing = data.find(entry => entry.missing)?.missing ?? [];
    const items = data.filter(entry => entry.item).map(entry => entry.item);

    renderMissingSets(missing);
    renderInventoryItems(items);
    statusText.textContent = `${missing.length} incomplete set${missing.length === 1 ? "" : "s"} found`;
  }
  catch (error) {
    summary.hidden = true;
    itemResults.hidden = true;
    result.className = "error";
    result.textContent = String(error);
    statusText.textContent = "The check could not be completed";
  }
  finally {
    button.disabled = false;
    button.textContent = "Check inventory";
  }
});

type Me = {
  username: string;
};

function showLoggedIn(me: Me) {
  currentUser.textContent = me.username;
  authStatus.hidden = false;
  authForm.hidden = true;
  authError.hidden = true;
  passwordInput.value = "";
}

function showLoggedOut() {
  authStatus.hidden = true;
  authForm.hidden = false;
}

async function submitCredentials(path: "/login" | "/register") {
  if (!authForm.reportValidity()) {
    return;
  }

  loginButton.disabled = true;
  registerButton.disabled = true;
  authError.hidden = true;

  try {
    const response = await fetch(path, {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({ username: usernameInput.value, password: passwordInput.value }),
      credentials: "same-origin"
    });

    if (!response.ok) {
      throw new Error((await response.text()).trim());
    }

    showLoggedIn(await response.json());
  }
  catch (error) {
    authError.textContent = error instanceof Error ? error.message : String(error);
    authError.hidden = false;
  }
  finally {
    loginButton.disabled = false;
    registerButton.disabled = false;
  }
}

async function loadCurrentUser() {
  try {
    const response = await fetch("/me", { credentials: "same-origin" });
    if (response.ok) {
      showLoggedIn(await response.json());
    } else {
      showLoggedOut();
    }
  }
  catch {
    showLoggedOut();
  }
}

authForm.addEventListener("submit", event => {
  event.preventDefault();
  submitCredentials("/login");
});

registerButton.addEventListener("click", () => {
  submitCredentials("/register");
});

logoutButton.addEventListener("click", async () => {
  logoutButton.disabled = true;
  try {
    await fetch("/logout", { method: "POST", credentials: "same-origin" });
  }
  finally {
    logoutButton.disabled = false;
    showLoggedOut();
  }
});

loadCurrentUser();
