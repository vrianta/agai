package agai

var _404__ []byte = []byte(`
<!DOCTYPE html>
<html lang="en">

<head>
    <title>404 - Page Not Found | WorkersHub</title>
    <meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1, user-scalable=0">
    <link href="/css/bootstrap/bootstrap.min.css" rel="stylesheet">
    <link rel="stylesheet" href="/css/font-awesome/all.min.css">
    <link href="/css/contructor.css" rel="stylesheet">
    <style>
        body {
            background-color: var(--wh-navy-deep);
            color: white;
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
            text-align: center;
            overflow: hidden;
        }

        .error-container {
            max-width: 600px;
            padding: 2rem;
        }

        .error-code {
            font-size: clamp(8rem, 20vw, 12rem);
            font-weight: 900;
            line-height: 1;
            color: rgba(255, 255, 255, 0.05);
            position: absolute;
            top: 50%;
            left: 50%;
            transform: translate(-50%, -50%);
            z-index: -1;
            user-select: none;
        }

        .construction-icon {
            font-size: 5rem;
            color: var(--wh-warning);
            margin-bottom: 1.5rem;
            display: inline-block;
            animation: bounce 2s infinite;
        }

        @keyframes bounce {

            0%,
            20%,
            50%,
            80%,
            100% {
                transform: translateY(0);
            }

            40% {
                transform: translateY(-20px);
            }

            60% {
                transform: translateY(-10px);
            }
        }

        .btn-home {
            background-color: var(--wh-warning);
            color: #000;
            padding: 1rem 2.5rem;
            font-weight: 800;
            border-radius: 100px;
            transition: all 0.3s ease;
            text-decoration: none;
            display: inline-block;
            margin-top: 2rem;
        }

        .btn-home:hover {
            transform: scale(1.05);
            background-color: #fff;
            color: var(--wh-navy-deep);
        }

        .tape-strip {
            height: 20px;
            width: 100%;
            background: repeating-linear-gradient(45deg,
                    var(--wh-warning),
                    var(--wh-warning) 20px,
                    #000 20px,
                    #000 40px);
            position: fixed;
            bottom: 0;
            left: 0;
            opacity: 0.6;
        }
    </style>
</head>

<body>

    <div class="error-code">404</div>

    <div class="error-container">
        <div class="construction-icon">
            <i class="fa-solid fa-conveyor-belt-arm"></i>
        </div>

        <h1 class="display-4 fw-bold mb-3">Wrong Turn at the <span class="text-warning">Site.</span></h1>

        <p class="lead text-white-50 mb-4">
            The page you are looking for has been moved, demolished, or never existed in the first place. Let's get you
            back to the main portal.
        </p>

        <div class="row justify-content-center g-3">
            <div class="col-12 col-md-8">
                <div class="d-flex flex-wrap gap-3 justify-content-center align-items-center">

                    <a href="/" class="btn-home shadow text-decoration-none">
                        <i class="fa-solid fa-house me-2"></i> Back to Home
                    </a>

                    <a href="#" onclick="history.back()"
                        class="btn btn-outline-light rounded-4 px-4 py-2 fw-bold border-2 d-flex align-items-center justify-content-center"
                        style="width: 180px; height: 52px;">
                        Go Back
                    </a>

                </div>
            </div>
        </div>
    </div>

    <div class="tape-strip"></div>

    <script src="/js/bootstrap/bootstrap.bundle.min.js"></script>
    <script src="/js/font-awesome/all.min.js"></script>
</body>

</html>
`)
