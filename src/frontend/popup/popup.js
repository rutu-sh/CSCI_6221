// Function to open subscription page
function openPage(provider) {
    // Open respective subscription page based on provider
    switch(provider) {
        case 'netflix':
            window.open('https://www.netflix.com', '_blank');
            break;
        case 'amazon':
            window.open('https://www.amazon.com/gp/video/storefront/', '_blank');
            break;
        default:
            var url = document.getElementById('subscription-url').value;
            console.log('Opening subscription page for:', url);
            if (url.startsWith('https') || url.startsWith('http') || url.startsWith('www')){
                window.open(url, '_blank');
            } else {
                console.error('Invalid provider URL:', provider);
                // TODO:: Show error message to user
            }
            break;
        }
}

// Function to open subscription settings page
function openSettings(provider) {
    console.log('Opening subscription settings page for:', provider);
    // Open respective subscription settings page based on provider
    switch(provider) {
        case 'netflix':
            window.open('https://www.netflix.com/YourAccount', '_blank');
            break;
        case 'amazon':
            window.open('https://www.amazon.com/mc/yourmembershipsandsubscriptions', '_blank');
            break;
        default:
            var url = document.getElementById('subscription-cancel-url').value;
            if (url.startsWith('https') || url.startsWith('http') || url.startsWith('www')){
                window.open(url, '_blank');
            } else {
                console.error('Invalid provider URL:', provider);
                // TODO:: Show error message to user
            }
            break;
        }
    }

// Function to handle click on "View" buttons
document.querySelectorAll('.view-btn').forEach(item => {
    item.addEventListener('click', event => {
        openPage(event.target.dataset.provider);
    });
});

// Function to handle click on "Cancel" buttons
document.querySelectorAll('.cancel-btn').forEach(item => {
    item.addEventListener('click', event => {
        openSettings(event.target.dataset.provider);
    });
});

// Function to handle click on "Add Subscription" button
document.getElementById('add-subscription-btn').addEventListener('click', function() {
    var addForm = document.getElementById('add-subscription-form');
    addForm.style.display = 'block';
});

// Function to handle click on "Save" button
document.getElementById('save-subscription-btn').addEventListener('click', function() {
    // Get input field values
    var name = document.getElementById('subscription-name').value;
    var url = document.getElementById('subscription-url').value;
    var cancelUrl = document.getElementById('subscription-cancel-url').value;
    var startDate = document.getElementById('subscription-start-date').value;

    // Create a new list item element
    var listItem = document.createElement('li');

    // Set the inner HTML of the list item with subscription details
    listItem.innerHTML = `
        <h3>${name}</h3>
        <button class="view-btn" id="subscription-url" data-provider=${url}>View</button>
        <button class="cancel-btn" id="subscription-cancel-url" data-provider=${cancelUrl}>Cancel</button>
    `;

    // Append the list item to the subscription list
    document.getElementById('subscription-list').appendChild(listItem);

    // Reset form and hide it
    document.getElementById('add-subscription-form').style.display = 'none';
    document.getElementById('subscription-url').value = '';
    document.getElementById('subscription-cancel-url').value = '';
    document.getElementById('subscription-start-date').value = '';
});
