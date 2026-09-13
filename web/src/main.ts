const input = document.getElementById("inventoryInput") as HTMLTextAreaElement;
const button = document.getElementById("submitButton") as HTMLButtonElement;
const result = document.getElementById("result") as HTMLPreElement;

button.addEventListener("click", async () => {
    const inventory = input.value.trim();

    try {
        JSON.parse(inventory);

        button.disabled = true;
        button.textContent = "Loading...";
        result.textContent = "Loading...";

        const response = await fetch("/inventory", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: inventory
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