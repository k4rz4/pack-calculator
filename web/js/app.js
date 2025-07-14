// Pack size management functions
function addPackSize() {
  const container = document.getElementById("pack-sizes");
  const newRow = document.createElement("div");
  newRow.className = "pack-size-row";
  newRow.innerHTML = `
        <input type="number" name="pack_sizes" class="pack-size-input" value="100" min="1" required>
        <button type="button" class="remove-btn" onclick="removePackSize(this)">Remove</button>
    `;
  container.appendChild(newRow);
  updateRemoveButtons();
}

function removePackSize(button) {
  const container = document.getElementById("pack-sizes");
  if (container.children.length > 1) {
    button.closest(".pack-size-row").remove();
    updateRemoveButtons();
  }
}

function updateRemoveButtons() {
  const container = document.getElementById("pack-sizes");
  const removeButtons = container.querySelectorAll(".remove-btn");
  removeButtons.forEach((btn, index) => {
    btn.disabled = container.children.length <= 1;
  });
}

function setPreset(sizes) {
  const container = document.getElementById("pack-sizes");
  container.innerHTML = "";

  sizes.forEach((size) => {
    const row = document.createElement("div");
    row.className = "pack-size-row";
    row.innerHTML = `
            <input type="number" name="pack_sizes" class="pack-size-input" value="${size}" min="1" required>
            <button type="button" class="remove-btn" onclick="removePackSize(this)">Remove</button>
        `;
    container.appendChild(row);
  });

  updateRemoveButtons();
}

// Handle successful response
document.addEventListener("htmx:beforeSwap", function (evt) {
  evt.preventDefault();
  try {
    const response = JSON.parse(evt.detail.xhr.responseText);

    if (response.success && response.data) {
      const data = response.data;
      let packsHtml = "";

      for (const [packSize, quantity] of Object.entries(data.packs_used)) {
        if (quantity > 0) {
          packsHtml += `
            <div class="result-item">
              <div class="result-number">${quantity}</div>
              <div class="result-label">Pack${quantity > 1 ? "s" : ""} of ${packSize} items</div>
            </div>
          `;
        }
      }

      const overage =
        data.items_overage > 0 ? ` (+${data.items_overage} extra)` : "";

      evt.detail.target.innerHTML = `
        <div class="result-container">
          <div class="result-title">✅ Optimal Pack Distribution</div>

          <div class="result-grid">
            ${packsHtml}
          </div>

          <div class="summary-grid">
            <div class="summary-item">
              <div class="summary-label">📦 Total Packs</div>
              <div class="summary-value">${data.total_packs}</div>
            </div>
            <div class="summary-item">
              <div class="summary-label">📊 Total Items</div>
              <div class="summary-value">${data.total_items}${overage}</div>
            </div>
            <div class="summary-item">
              <div class="summary-label">⚡ Time</div>
              <div class="summary-value">${data.calculation_time}</div>
            </div>
          </div>
        </div>
      `;
    } else {
      throw new Error(response.error || "Invalid response format");
    }
  } catch (error) {
    evt.detail.target.innerHTML = `
      <div class="error-container">
        ❌ Error: ${error.message}
      </div>
    `;
  }
});

// Handle HTTP errors
document.addEventListener("htmx:responseError", function (evt) {
  evt.detail.target.innerHTML = `
        <div class="error-container">
            ❌ Error: Failed to connect to the API
        </div>
    `;
});

// Initialize on page load
document.addEventListener("DOMContentLoaded", function () {
  updateRemoveButtons();
});
