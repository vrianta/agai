<!DOCTYPE html>
<html data-bs-theme="dark" lang="en">

<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>First User Registration - OP Resume</title>
  <link href="/css/Bootstrap/bootstrap.min.css" rel="stylesheet">

  <style>
    body {
      background-color: var(--bs-body-bg);
      display: flex;
      align-items: center;
      justify-content: center;
      min-height: 100vh;
      padding: 1rem;
    }
    .registration-card {
      max-width: 480px;
      width: 100%;
      border-radius: 0.75rem;
      box-shadow: 0 0.5rem 1rem rgba(0,0,0,.15);
      animation: fadeIn 0.6s ease;
    }
    @keyframes fadeIn {
      from { opacity: 0; transform: translateY(10px); }
      to { opacity: 1; transform: translateY(0); }
    }
    .logo-placeholder {
      width: 80px;
      height: 80px;
      border-radius: 50%;
      background-color: rgba(255,255,255,0.05);
      display: flex;
      align-items: center;
      justify-content: center;
      margin: 0 auto 1.5rem auto;
      color: var(--bs-secondary-color);
      font-size: 0.8rem;
      text-align: center;
    }
    footer {
      margin-top: 2rem;
      font-size: 0.875rem;
      text-align: center;
      color: var(--bs-secondary-color);
      max-width: 480px;
    }
  </style>
</head>

<body>
  <main>
    <!-- Registration Card -->
    <div class="card registration-card bg-body-tertiary p-4 p-md-5">
      
      <!-- Round Logo Placeholder -->
      <div class="logo-placeholder">
        Logo
      </div>

      <!-- Title -->
      <h4 class="mb-4 text-center">Create Your First OP Resume Account</h4>

      <!-- Error Section -->
      <div id="errorAlert" class="alert alert-danger d-none" role="alert">
        <!-- Example: "Email already exists" -->
      </div>

      <!-- Registration Form -->
      <form method="post" novalidate>
        
        <!-- First Name -->
        <div class="form-floating mb-3">
          <input type="text" class="form-control" id="firstName" name="firstName" placeholder="First Name" required>
          <label for="firstName">First Name</label>
        </div>

        <!-- Last Name -->
        <div class="form-floating mb-3">
          <input type="text" class="form-control" id="lastName" name="lastName" placeholder="Last Name" required>
          <label for="lastName">Last Name</label>
        </div>

        <!-- Email -->
        <div class="form-floating mb-3">
          <input type="email" class="form-control" id="email" name="email" placeholder="Email Address" required>
          <label for="email">Email Address</label>
        </div>

        <!-- Password -->
        <div class="form-floating mb-3">
          <input type="password" class="form-control" id="password" name="password" placeholder="Password" required>
          <label for="password">Password</label>
        </div>

        <!-- Confirm Password -->
        <div class="form-floating mb-4">
          <input type="password" class="form-control" id="confirmPassword" name="confirmPassword" placeholder="Confirm Password" required>
          <label for="confirmPassword">Confirm Password</label>
        </div>

        <!-- Submit -->
        <div class="d-grid">
          <button type="submit" class="btn btn-primary btn-lg">Create Account</button>
        </div>
      </form>

      <?php if ($$error): ?>
        <div class="text-danger">
          Error : <?= $$error ?> 
        </div>
        <?php endif; ?>
    </div>



    <!-- Footer -->
    <footer>
      <p>
        <strong>OP Resume</strong> is an open-source project that helps you create and host your professional resume online.  
        Deploy it on your own server, customize it, and maintain full control of your data — free for personal and non-commercial use.
      </p>
    </footer>
  </main>

  <script src="/Js/Bootstrap/bootstrap.min.js"></script>

  <!-- Example JS to Show Error -->
  <script>
    // Example usage:
    // document.getElementById("errorAlert").textContent = "Email already exists";
    // document.getElementById("errorAlert").classList.remove("d-none");
  </script>
</body>

</html>
