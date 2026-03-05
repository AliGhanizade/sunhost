function togglePassword(inputId, iconElement) {
  const input = document.getElementById(inputId);
  if (input.type === "password") {
    input.type = "text";
    iconElement.className = "bi bi-eye-slash password-toggle-right";
  } else {
    input.type = "password";
    iconElement.className = "bi bi-eye password-toggle-right";
  }
}

function showError(elementId, message) {
  const errorElement = document.getElementById(elementId);
  errorElement.innerHTML = message;
  errorElement.style.display = "block";
}

function hideError(elementId) {
  document.getElementById(elementId).style.display = "none";
}

async function login(event) {
  event.preventDefault();
  hideError("error-login");

  const loginData = {
    username: document.getElementById("username-login").value,
    password: document.getElementById("password-login").value,
  };

  try {
    const response = await fetch("/api/users/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(loginData),
    });

    const result = await response.json();

    if (response.ok) {
      localStorage.setItem("isLoggedIn", "true");
      localStorage.setItem("username", loginData.username);
      localStorage.setItem("fullName", result.user.full_name);
      localStorage.setItem("email", result.user.email);
      window.location.href = "/panel";
    } else {
      showError(
        "error-login",
        result.error || "نام کاربری یا رمز عبور اشتباه است",
      );
    }
  } catch (err) {
    showError(
      "error-login",
      result.error || "نام کاربری یا رمز عبور اشتباه است",
    );
  }
}

async function register(event) {
  event.preventDefault();
  hideError("error-register");
  const username = document.getElementById("username-register").value;
  const firstName = document.getElementById("first-name").value;
  const lastName = document.getElementById("last-name").value;
  const password = document.getElementById("password-register").value;
  const confirmPassword = document.getElementById(
    "password-confirm-register",
  ).value;

  const userData = {
    full_name: (firstName + " " + lastName).trim(),
    username: document.getElementById("username-register").value,
    email: document.getElementById("email-register").value,
    password: password,
  };

  if (username.length < 4) {
    showError("error-register", "نام کاربری باید حداقل 4 کاراکتر باشد.");
    return;
  }

  if (lastName.length < 3) {
    showError("error-register", "نام خانوادگی باید حداقل 3 کاراکتر باشد.");
    return;
  }

  if (password.length < 4) {
    showError("error-register", "رمز عبور باید حداقل 4 کاراکتر باشد.");
    return;
  }

  if (password !== confirmPassword) {
    showError("error-register", "رمز عبور و تکرار آن با هم مطابقت ندارند.");
    return;
  }

  try {
    const response = await fetch("/api/users/register", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(userData),
    });

    const result = await response.json();

    if (response.ok) {
      alert("ثبت‌نام موفقیت‌آمیز بود!");
      window.location.reload();
    } else {
      showError("error-register", result.error || "خطایی در ثبت نام رخ داد");
    }
  } catch (err) {
    showError("error-register", "ارتباط با سرور برقرار نشد");
  }
}
