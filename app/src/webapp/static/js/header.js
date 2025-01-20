// document.addEventListener('DOMContentLoaded', function() {
//     const OrganizationDropdownButton = document.getElementById('OrganizationDropdownButton');
//     const OrganizationItems = document.getElementById('OrganizationItems');

//     OrganizationDropdownButton.addEventListener('click', function(event) {
//         OrganizationItems.classList.toggle('hidden');
//         event.stopPropagation(); // Prevent the click event from propagating to the document
//     });

//     // Close dropdown if clicked outside
//     document.addEventListener('click', function(event) {
//         if (!OrganizationDropdownButton.contains(event.target) && !OrganizationItems.contains(event.target)) {
//             OrganizationItems.classList.add('hidden');
//         }
//     });

//     // Prevent the click event from propagating to the document
//     OrganizationItems.addEventListener('click', function(event) {
//         event.stopPropagation();
//     });
// });
// document.addEventListener('DOMContentLoaded', function() {
//     document.getElementById('OrganizationDropdownButton').addEventListener('click', function(event) {
//         console.log('clicked');
//         console.log
//         event.stopPropagation(); // Prevent the click event from bubbling up to the document
//         var OrganizationItems = document.getElementById('OrganizationItems');
//         if (OrganizationItems.classList.contains('hidden')) {
//             OrganizationItems.classList.remove('hidden');
//         } else {
//             OrganizationItems.classList.add('hidden');
//         }
//     });

//     // Close the dropdown if clicked outside
//     document.addEventListener('click', function(event) {
//         var OrganizationItems = document.getElementById('OrganizationItems');
//         if (!document.getElementById('OrganizationDropdownButton').contains(event.target)) {
//             OrganizationItems.classList.add('hidden');
//         }
//     });
// });