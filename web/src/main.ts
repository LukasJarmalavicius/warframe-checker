const input = document.getElementById("inventoryInput") as HTMLTextAreaElement;
const button = document.getElementById("submitButton") as HTMLButtonElement;
const result = document.getElementById("result") as HTMLPreElement;
const test = document.getElementById("test") as HTMLPreElement;
const check = document.getElementById("jsonCheck") as HTMLInputElement;

button.addEventListener("click", async () => {
  const inventory = input.value.trim()

  let json
  if (!check.checked) {
    const items = inventory.split("\n").map(line => line.trim()).filter(line => line.length > 0).map(name => ({ name, quantity: 1 }))
    json = JSON.stringify({ items })
  } else {
    json = inventory
  }

  try {
    test.textContent = (check.checked ? "json\n" : "nojson\n") + json;

    button.disabled = true;
    button.textContent = "Loading...";
    result.textContent = "Loading...";

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

    result.textContent = JSON.stringify(data, null, 2);
  }
  catch (error) {
    result.textContent = String(error);
  }
  finally {
    button.disabled = false;
    button.textContent = "Submit";
  }
});
