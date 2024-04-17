//-------------------------------------------------------------
// Manage subscriptions page
//-------------------------------------------------------------

document.addEventListener("DOMContentLoaded", function () {
  const subscriptionsContainer = document.getElementById("subscriptions");
  const noCardsContainer = document.getElementById("no-cards");
  const subscriptionsHeading = document.getElementById("subscriptions-heading");

  const username = sessionStorage.getItem("username");
  if (!username) {
    console.error("Username not found in sessionStorage.");
    alert("Username not found.");
    return;
  }
  
  // API call to get subscriptions data
  fetch(`https://7se83qeyid.execute-api.us-east-1.amazonaws.com/dev/v2/subscriptions?username=${username}`)
    .then((response) => {
      if (!response.ok) {
        throw new Error("Network response was not ok");
      }
      return response.json();
    })
    .then((data) => {
      console.log(data); // Verify that the data is retrieved correctly
      // Check if subscriptions data exists and is an array
      if (!data || data.length === 0) {
        noCardsContainer.innerHTML = "¯\\_(ツ)_/¯ No subscriptions to show"; // Display the shrug emoji
        subscriptionsHeading.style.display = "none"; // Hide the heading
      } else {
        allSubscriptions = localStorage.setItem("allSubscriptions", JSON.stringify(data));
        // Iterate over each subscription object and create subscription cards
        data.forEach((subscription) => {
          console.log("subzzzzz",subscription); // Verify that the subscription data is correct
          const name = subscription.name;
          const url = subscription.url;
          const settingsUrl = subscription.settings_url;
          const plan = subscription.plan;
          const startDate = subscription.start_date;
          const cost = subscription.cost;
          let icon = subscription.icon;
          const lastPaymentDate = subscription.last_payment_date;
          const category = subscription.category;
          const subscriptionId = subscription.uuid;

          // Check if the URL contains "netflix"
        if (url.includes("netflix")) {
          // Set the icon to the Netflix logo
          icon = "./logos/netflix.png";
        }
        if (url.includes("spotify")) {
          // Set the icon to the Spotify logo
          icon = "./logos/spotify.png";
        }
        if (url.includes("hulu")) {
          // Set the icon to the Amazon logo
          icon = "./logos/hulu.png";
        }
        if (url.includes("instacart")) {
          // Set the icon to the Instacart logo
          icon = "./logos/instacart.png";
        }
        if (url.includes("mint")) {
          // Set the icon to the Mint logo
          icon = "./logos/mint.png";
        }

          // Now you can use these variables to create your subscription card or perform other operations
          console.log("Name:", name);
          console.log("URL:", url);
          console.log("Settings URL:", settingsUrl);
          console.log("Plan:", plan);
          console.log("Start Date:", startDate);
          console.log("Cost:", cost);
          console.log("Icon:", icon);
          console.log("Last Payment Date:", lastPaymentDate);
          console.log("Category:", category);
          console.log("Subscription ID:", subscriptionId);

          // Create subscription card using these variables
          const subscriptionCard = generateSubscriptionCard({
            name: name,
            url: url,
            settingsUrl: settingsUrl,
            plan: plan,
            start_date: startDate,
            cost: cost,
            icon: icon,
            last_payment_date: lastPaymentDate,
            category: category,
            subscriptionId: subscriptionId
          });

          // Append subscription card to container
          subscriptionsContainer.appendChild(subscriptionCard);
        });
      }
    })
    .catch((error) => {
      console.error("There was a problem fetching subscription data:", error);
    });
});


function generateSubscriptionCard(subscription) {
  const card = document.createElement("div");
  card.classList.add("card");

  const icon = document.createElement("img");
  icon.classList.add("subscription-icon");
  icon.src = subscription.icon;

  const detailsContainer = document.createElement("div");
  detailsContainer.classList.add("subscription-details-container");

  const nameAndPlanContainer = document.createElement("div");
  nameAndPlanContainer.classList.add("name-and-plan-container");

  const name = document.createElement("div");
  name.classList.add("subscription-name");
  name.textContent = subscription.name;

  const plan = document.createElement("div");
  plan.classList.add("subscription-plan");
  plan.textContent = `${subscription.plan}`;

  nameAndPlanContainer.appendChild(name);
  nameAndPlanContainer.appendChild(plan);

  const cost = document.createElement("div");
  cost.classList.add("subscription-cost");
  cost.textContent = `$${subscription.cost}`;

  const nextPaymentDate = document.createElement("div");
  nextPaymentDate.classList.add("subscription-next-payment-date");
  nextPaymentDate.textContent = `${subscription.last_payment_date}`;

  const nextPaymentLabel = document.createElement("div");
  nextPaymentLabel.classList.add("next-payment-label");
  nextPaymentLabel.textContent = "Last Payment Date";

  const settingsBtn = document.createElement("button");
  settingsBtn.classList.add("settings-btn");
  settingsBtn.textContent = "";
  settingsBtn.addEventListener("click", function (event) {
    document.querySelectorAll(".settings-dropdown").forEach((dropdown) => {
      dropdown.remove();
    });
    localStorage.setItem("subscription", JSON.stringify(subscription));
    const dropdown = document.createElement("div");
    dropdown.classList.add("settings-dropdown");
    dropdown.innerHTML = `
      <a href="update.html"><img src="update.png" alt="Update" />Update</a>
      <a href="${subscription.settingsUrl}" target="_blank"><img src="deets.png" alt="Details" />Details</a>
      <a href="#" id="delete-link" class="delete-link"><img src="delete.png" alt="Delete" />Delete</a>
    `;

    const cardRect = card.getBoundingClientRect();
    dropdown.style.top = `${cardRect.bottom - 30}px`;
    document.body.appendChild(dropdown);

    // Close the dropdown when clicking outside of it or on the settings button
    document.addEventListener("click", function (e) {
      if (!dropdown.contains(e.target) && e.target !== settingsBtn) {
        dropdown.remove();
      }
    });

    // Prevent the default action of the delete link and show confirmation
    const deleteLink = dropdown.querySelector("#delete-link");

    deleteLink.addEventListener("click", function (event) {
      event.preventDefault();
      const subscriptionId = subscription.subscriptionId; // Retrieve subscriptionId
      confirmDelete(subscriptionId);
      event.stopPropagation();
    });
  });

  // Define the confirmDelete function
  function confirmDelete(subscriptionId) {
    console.log("Subscription ID to delete:", subscriptionId);
    if (confirm("Are you sure you want to delete the subscription?")) {
      const username = sessionStorage.getItem("username");
      if (!username) {
        console.error("Username not found in sessionStorage.");
        alert("Username not found.");
        return;
      }

      fetch(`https://7se83qeyid.execute-api.us-east-1.amazonaws.com/dev/v2/subscriptions/${subscriptionId}?username=${username}`, {
        method: "DELETE"
      })
        .then(response => {
          if (!response.ok) {
            throw new Error("Failed to delete subscription");
          }
          // remove the subscription card from the UI
          card.remove();
        })
        .catch(error => {
          console.error("Error deleting subscription:", error);
        });
    }
  }


  detailsContainer.appendChild(icon);
  detailsContainer.appendChild(nameAndPlanContainer);
  detailsContainer.appendChild(cost);
  const dateContainer = document.createElement("div");
  dateContainer.classList.add("subscription-date-container");
  dateContainer.appendChild(nextPaymentDate);
  dateContainer.appendChild(nextPaymentLabel);
  detailsContainer.appendChild(dateContainer);
  detailsContainer.appendChild(settingsBtn);

  const pastPaymentsContainer = document.createElement("div");
  pastPaymentsContainer.classList.add("past-payments-container");

  const pastPaymentsList = document.createElement("ul");
  pastPaymentsList.classList.add("past-payments-list");

  const title = document.createElement("h3");
  title.textContent = "Payment History";
  pastPaymentsList.appendChild(title);

  const pastPayment = document.createElement("li");
  pastPayment.classList.add("past-payment");
  pastPayment.textContent = subscription.last_payment_date; // Set the text content to the last payment date
  pastPaymentsList.appendChild(pastPayment); // Append the list item to the pastPaymentsList
  pastPaymentsContainer.appendChild(pastPaymentsList); // Append the list to the pastPaymentsContainer
  

  card.appendChild(detailsContainer);
  card.appendChild(pastPaymentsContainer);

  // Toggle function for past payments accordion
  function togglePastPayments() {
    const pastPaymentsList = card.querySelector(".past-payments-list");
    pastPaymentsList.classList.toggle("show");
  }

  // Add click event listener to card for toggling past payments
  card.addEventListener("click", togglePastPayments);

  return card;
}


const darkModeToggle = document.getElementById("mode-toggle");
const body = document.body;

darkModeToggle.addEventListener("change", () => {
  body.classList.toggle("light-mode");
  const settingsBtns = document.querySelectorAll(".settings-btn");
  settingsBtns.forEach((btn) => {
    btn.classList.toggle("light-mode");
  });
});

//-------------------------------------------------------------
// Login form
//-------------------------------------------------------------

document.addEventListener("DOMContentLoaded", function () {
  const loginButton = document.querySelector(".login-button");
  const popupContainer = document.createElement("div");
  popupContainer.classList.add("popup-container");

  const loginForm = document.createElement("form");
  loginForm.classList.add("login-form");

  const usernameInput = document.createElement("input");
  usernameInput.setAttribute("type", "text");
  usernameInput.setAttribute("placeholder", "Username");

  const passwordInput = document.createElement("input");
  passwordInput.setAttribute("type", "password");
  passwordInput.setAttribute("placeholder", "Password");

  const newHere = document.createElement("a");
  newHere.classList.add("new-here");
  newHere.textContent = "New here? Register";
  newHere.setAttribute("href", "./registration.html");

  const forgotPassword = document.createElement("a");
  forgotPassword.classList.add("forgot-password");
  forgotPassword.textContent = "Forgot password?";
  forgotPassword.setAttribute("href", "./forgot-password.html");

// Add event listener to handle password reset process
forgotPassword.addEventListener("click", function(event) {
  event.preventDefault(); // Prevent default link behavior

  // Retrieve the username from sessionStorage
  const username = sessionStorage.getItem("username");

  // Send a POST request to the forgot password API endpoint
  fetch("https://fw79wa9t2e.execute-api.us-east-1.amazonaws.com/dev/auth-forgot-password", {
      method: "POST",
      headers: {
          "Content-Type": "application/json"
      },
      body: JSON.stringify({ username: username })
  })
  .then(response => {
      if (!response.ok) {
          throw new Error("Failed to initiate password reset process");
      }
      // Handle successful response (e.g., show a success message)
      alert("OTP sent successfully. Please check your email.");
      window.location.href = "./forgot-password.html";
  })
  .catch(error => {
      console.error("Error:", error);
      // Handle error (e.g., show an error message)
      alert("Failed to initiate password reset process. Please try again later.");
  });
});

// Append the 'Forgot password?' link to the document body or appropriate container
document.body.appendChild(forgotPassword);

  const submitButton = document.createElement("button");
  submitButton.setAttribute("type", "submit");
  submitButton.textContent = "Login";

  loginForm.appendChild(usernameInput);
  loginForm.appendChild(passwordInput);
  loginForm.appendChild(newHere);
  loginForm.appendChild(forgotPassword);
  loginForm.appendChild(submitButton);

  popupContainer.appendChild(loginForm);

  function togglePopup() {
    popupContainer.classList.toggle("active");
  }

  loginButton.addEventListener("click", togglePopup);

  // Close popup when clicking outside of it
  document.addEventListener("click", function (event) {
    const isClickInsidePopup = popupContainer.contains(event.target);
    const isClickOnLoginButton = loginButton.contains(event.target);
    if (!isClickInsidePopup && !isClickOnLoginButton) {
      popupContainer.classList.remove("active");
    }
  });

  document.body.appendChild(popupContainer);
  // Add event listener to the login form submission
  loginForm.addEventListener("submit", function (event) {
    event.preventDefault(); // Prevent the default form submission

    // Get the username and password from the form inputs
    const username = usernameInput.value;
    const password = passwordInput.value;

    // Make the API call
    var loggedInUser = username;
    console.log("logging in"); // Verify that loggedInUser is undefined
    fetch(
      "https://fw79wa9t2e.execute-api.us-east-1.amazonaws.com/dev/auth-signin",
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          username: username,
          password: password,
        }),
      }
    )
      .then((response) => {
        if (!response.ok) {
          throw new Error("Network response was not ok");
        }
        return response.json();
      })
      .then((data) => {
        console.log(data); // Verify that the username is retrieved correctly
        // Store the username in sessionStorage
        sessionStorage.setItem("username", username);
        // Set the loggedInUser variable
        loggedInUser = username;
        console.log(loggedInUser); // Verify that loggedInUser holds the correct value

        // Display the loggedInUser in the HTML
        const loggedInUserSpan = document.getElementById("loggedInUser");
        console.log(loggedInUserSpan); // Verify that the element is found
        if (loggedInUserSpan) {
          loggedInUserSpan.textContent = loggedInUser;
        } else {
          console.error("Element with id 'loggedInUser' not found.");
        }
        // If login is successful, redirect to manage.html
        window.location.href = "manage.html";
      })
      .catch((error) => {
        console.error("There was a problem with the fetch operation:", error);
      });
    console.log("logged in user: " + loggedInUser); // Verify that loggedInUser is undefined
  });
});

document.addEventListener("DOMContentLoaded", function () {
  const loginButton = document.querySelector(".new-here");

  loginButton.addEventListener("click", function () {
    window.location.href = "registration.html";
  });
});

// function confirmDelete(subscriptionId) {
//   console.log("Subscription ID in delete:", subscriptionId);
//   if (confirm("Are you sure you want to delete the subscription?")) {
//     // get the subscription ID from the card
//     const username = sessionStorage.getItem("username");
//     if (!username) {
//       console.error("Username not found in sessionStorage.");
//       alert("Username not found.");
//       return;
//     }

//     fetch(`https://7se83qeyid.execute-api.us-east-1.amazonaws.com/dev/v2/subscriptions/${subscriptionId}?username=${username}`, {
//       method: "DELETE"
//     })
//     .then(response => {
//       if (!response.ok) {
//         throw new Error("Failed to delete subscription");
//       }
//       // Optionally, you can remove the subscription card from the UI here
//       // For example: subscriptionCard.remove();
//     })
//     .catch(error => {
//       console.error("Error deleting subscription:", error);
//     });
//   }
// }



//-------------------------------------------------------------
// Registration form
//-------------------------------------------------------------
const registrationForm = document.getElementById("registration-form");

registrationForm.addEventListener("submit", function (event) {
  event.preventDefault(); // Prevent the default form submission

  // Get the form inputs
  const name = document.getElementById("name").value;
  const email = document.getElementById("email").value;
  const password = document.getElementById("password").value;
  const confirmPassword = document.getElementById("confirm-password").value;
  const username = document.getElementById("username").value;

  // Make sure password and confirm password match
  if (password !== confirmPassword) {
    alert("Passwords do not match!");
    return;
  }

  // Make the API call for registration
  fetch(
    "https://fw79wa9t2e.execute-api.us-east-1.amazonaws.com/dev/auth-signup",
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        email: email,
        password: password,
        name: name,
        username: username,
      }),
    }
  )
    .then((response) => {
      if (!response.ok) {
        throw new Error("Network response was not ok");
      }
      return response.json();
    })
    .then((data) => {
      // If registration is successful, get the username from the response
      console.log("Registration successful. Username:", username);

      // store the username in sessionStorage
      sessionStorage.setItem("username", username);
      console.log("username in session", username);

      // Redirect user to enter OTP page
      window.location.href = "enter-otp.html";

      // Call function to submit OTP with retrieved username
      submitOTP(username);
    })
    .catch((error) => {
      console.error("There was a problem with the fetch operation:", error);
    });
});
