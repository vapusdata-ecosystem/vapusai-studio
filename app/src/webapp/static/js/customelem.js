function setupDropdown({
    dropdownToggleSelector,
    dropdownMenuSelector,
    inputFieldSelector = null, // Optional, for updating a hidden input or text field
    parentClass = 'parent',
    childClass = 'child',
    displayAttribute,
    valueTransform, // Default transform for display text
  }) {
    console.log(dropdownToggleSelector);
    const dropdownToggle = document.querySelector(dropdownToggleSelector);
    const dropdownMenu = document.querySelector(dropdownMenuSelector);
    const inputField = inputFieldSelector ? document.querySelector(inputFieldSelector) : null;
  
    // Toggle dropdown visibility
    dropdownToggle.addEventListener('click', () => {
      dropdownMenu.classList.toggle('show');
    });
  
    // Close dropdown when clicking outside
    document.addEventListener('click', (event) => {
      if (!dropdownToggle.contains(event.target) && !dropdownMenu.contains(event.target)) {
        dropdownMenu.classList.remove('show');
      }
    });
  
    // Handle selection of dropdown items
    dropdownMenu.addEventListener('click', (event) => {
      const item = event.target;
      // Ignore clicks on parent items
      if (item.classList.contains(parentClass)) {
        return;
      }
      // Handle selection of child items
      if (item.classList.contains(childClass)) {
        let displayValue = item.getAttribute(displayAttribute);
        // Update the toggle text
        if (valueTransform !== undefined) {
          displayValue = valueTransform(displayValue);
        }
        dropdownToggle.textContent = displayValue;
        // Update the input field, if provided
        if (inputField) {
          const selectedValue = item.dataset.value;
          inputField.value = selectedValue;
        }
        // Close the dropdown
        dropdownMenu.classList.remove('show');
      }
    });
  }
  


function addTablePagination({
  tableId,
  rowsPerPage,
  prevPageBtn,
  nextPageBtn,
  currentPageSizeSpan,
  totalElementsSpan,
}) {
  console.log(tableId,"--00");
  const table = document.getElementById(tableId);
  const rows = table.querySelectorAll("tbody tr");
  console.log(rows,"--1");
  const totalRows = rows.length;
  console.log(totalRows,"--2");
  const totalPages = Math.ceil(totalRows / rowsPerPage);
  let currentPage = 1;

  // Function to display rows for the current page
  function displayPage(page) {
    console.log(page,"--3");
    const start = (page - 1) * rowsPerPage;
    const end = Math.min(page * rowsPerPage, totalRows);
    console.log(start, end);
    rows.forEach((row, index) => {
      row.style.display = index >= start && index < end ? "" : "none";
    });

    // Update pagination info
    currentPageSizeSpan.textContent = end - start;
    totalElementsSpan.textContent = totalRows;

    // Enable/Disable navigation buttons
    prevPageBtn.disabled = page === 1;
    nextPageBtn.disabled = page === totalPages;
  }

  // Add event listeners for navigation
  prevPageBtn.addEventListener("click", () => {
    if (currentPage > 1) {
      currentPage--;
      displayPage(currentPage);
    }
  });

  nextPageBtn.addEventListener("click", () => {
    console.log(currentPage);
    if (currentPage < totalPages) {
      currentPage++;
      displayPage(currentPage);
    }
  });

  // Initialize table display
  displayPage(currentPage);
}
