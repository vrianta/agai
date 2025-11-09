<!DOCTYPE html>
<html data-bs-theme="dark" lang="en">

<head>
    <title>AVS CMS</title>
    <link href="/css/Bootstrap/bootstrap.min.css" rel="stylesheet">
    <link href="/css/main.css" rel="stylesheet">
</head>

<body>
    <main>
        <section class="vh-100 bg-body">
            <div class="container h-100">
                <div class="row align-items-center h-100">
                    <!-- Left column with text -->
                    <div class="col-lg-6 mb-5 mb-lg-0 text-center text-lg-start">
                        <h1 class="my-5 display-4 fw-bold text-primary">
                            Welcome to AVS<br />
                            <span class="text-secondary fs-1">Content Management System</span>
                        </h1>
                        <p class="text-warning">
                            You have accessed a system managed by Accenture. You are required to have authorization from Accenture before you proceed and you are strictly limited to use set out within that authorization. Unauthorized access to or misuse of this system is prohibited and constitutes an offense under the Computer Misuse Act 1990. If you disclose any information obtained through this system without authority Accenture may take legal action against you.
                        </p>
                    </div>

                    <!-- Right column with login form -->
                    <div class="col-lg-6 d-flex justify-content-center">
                        <div class="card shadow-sm w-100" style="max-width: 400px;">
                            <div class="card-body p-5 bg-body-tertiary">
                                <h4 class="mb-4 text-center">Login</h4>
                                <form method="post">
                                    <!-- Email input -->
                                    <div class="form-floating mb-3">
                                        <input type="text" class="form-control" id="loginEmail"
                                            placeholder="name" name="loginEmail" value="<?= $$UserName ?>">
                                        <label for="loginEmail">User Name</label>
                                    </div>

                                    <!-- Password input -->
                                    <div class="form-floating mb-4">
                                        <input type="password" class="form-control" id="loginPassword"
                                            placeholder="Password" name="loginPassword" value="<?= $$Password ?>">
                                        <label for="loginPassword">Password</label>
                                    </div>

                                    <!-- Remember me checkbox -->
                                    <div class="form-check mb-4">
                                        <input class="form-check-input" type="checkbox" id="rememberMe" name="rememberMe">
                                        <label class="form-check-label" for="rememberMe">
                                            Remember me
                                        </label>
                                    </div>

                                    <!-- Submit button -->
                                    <div class="d-grid mb-3">
                                        <button type="submit" class="btn btn-primary">Login</button>
                                    </div>

                                    <!-- Forgot password and sign-up link -->
                                    <div class="text-center">
                                        <a href="#">Forgot password?</a>
                                    </div>
                                    <!-- Register User -->
                                    <div class="text-center">
                                        <a href="/register">Register User</a>
                                    </div>
                                </form>

                                <?php if ($$error): ?>
                                    <div class="text-danger">
                                        Error : <?= $$error ?>
                                    </div>
                                <?php endif; ?>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </section>

    </main>
    <footer>
        <script src="/js/Bootstrap/bootstrap.min.js"></script>
    </footer>
</body>

</html>