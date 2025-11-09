<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <title>
        <?= $$Title ?> | Portfolio
    </title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link rel="stylesheet" href="/css/Bootstrap/bootstrap.min.css">
    <link rel="stylesheet" href="/css/Bootstrap-Icons/bootstrap-icons.css">
    <link rel="stylesheet" href="/css/main.css">
</head>

<body data-bs-theme="dark" data-bs-spy="scroll" data-bs-target=".navbar" data-bs-offset="50">
    <?= include("components.header") /* Hello */ ?>
    <?= print("Testing print") ?>
    <main>
        <section id="home" class="hero-section text-white d-flex align-items-center" style="background-image: url(/static/img/hero-background.avif);">
            <div class="overlay"></div>
            <div class="container text-center position-relative z-1">
                <h1 class="display-4 fw-bold"><?= $$Hero->Heading ?></h1>
                <p class="lead mb-4"><?= $$Hero->SubHeading ?></p>
                <?php foreach ($$Hero->CallToActions as $key => $CallToAction): ?>
                    <a href="<?= $CallToAction->Href ?>" class="btn btn-primary me-2"><?= $CallToAction->Text ?></a>
                <?php endforeach ?>
                <!-- <a href="#protects" class="btn btn-outline-light">View Projects</a>
                <a href="#contact-me" class="btn btn-outline-light">Get in Touch</a> -->
            </div>
        </section>

        <section id="about-me" class="about-me py-5">
            <div class="container">
                <div class="row align-items-center">
                    <!-- Left: Image -->
                    <div class="col-md-5 text-center mb-4 mb-md-0">
                        <img src="<?= $$AboutMe->Picture ?>" alt="Your Photo" class="img-fluid rounded shadow">
                    </div>

                    <!-- Right: Details -->
                    <div class="col-md-7">
                        <h2 class="mb-3">About Me</h2>
                        <p class="lead"><?= $$AboutMe->AboutMe ?></p>
                        <ul class="list-unstyled mt-3">
                            <li><strong>Email:</strong> <?= $$ContactDetails->Email ?></li>
                            <li><strong>Location:</strong> <?= $$ContactDetails->Location ?></li>
                            <li><strong>Phone:</strong> <?= $$ContactDetails->Phone ?></li>
                        </ul>
                        <a href="#skills" class="btn btn-primary mt-3">Skills</a>
                    </div>
                </div>
            </div>
        </section>

        <section id="skills" class="skills-section py-5 bg-body-secondary">
            <div class="container">
                <h2 class="text-center mb-5">My Skills</h2>

                <div class="row">
                    <?php foreach ($$Skills as $key => $skill): ?>
                        <div class="col-md-6 mb-4">
                            <h5><?= $skill->Name ?></h5>
                            <div class="progress bg-body" style="height: 20px;">
                                <?php if ($skill->Level >= 90): ?>
                                    <div class="progress-bar bg-success" role="progressbar" style="width: <?= $skill->Level ?>%;"
                                        aria-valuenow="<?= $skill->Level ?>" aria-valuemin="0" aria-valuemax="100">
                                    <?php elseif ($skill->Level >= 85): ?>
                                        <div class="progress-bar bg-warning" role="progressbar"
                                            style="width: <?= $skill->Level ?>%;" aria-valuenow="<?= $skill->Level ?>"
                                            aria-valuemin="0" aria-valuemax="100">
                                        <?php elseif ($skill->Level >= 80): ?>
                                            <div class="progress-bar bg-primary" role="progressbar"
                                                style="width: <?= $skill->Level ?>%;" aria-valuenow="<?= $skill->Level ?>"
                                                aria-valuemin="0" aria-valuemax="100">
                                            <?php elseif ($skill->Level >= 75): ?>
                                                <div class="progress-bar bg-danger" role="progressbar" style="width: <?= $skill->Level ?>%;"
                                                    aria-valuenow="<?= $skill->Level ?>" aria-valuemin="0" aria-valuemax="100">
                                                <?php endif ?>

                                                <?= $skill->Level ?>%</div>
                                            </div>
                                        </div>
                                    <?php endforeach ?>
                                    </div>
                            </div>
        </section>

        <section id="experience" class="experience-section py-5 bg-body text-light">
            <div class="container">
                <h2 class="text-center mb-5">Experience</h2>

                <div class="timeline">
                    <?php foreach ($$Experiences as $key => $Experience): ?>
                        <div class="timeline-item mb-5 position-relative ps-4 border-start border-4 border-primary">
                            <span class="position-absolute top-0 start-0 translate-middle bg-primary rounded-circle"
                                style="width: 20px; height: 20px; border: 3px solid #212529;"></span>
                            <h5 class="mb-1"><?= $Experience->Designation ?></h5>
                            <small class="text-muted text-info"><?= $Experience->Duration ?></small>
                            <p class="mt-2 text-light"><?= $Experience->Description ?></p>
                        </div>
                    <?php endforeach; ?>
                </div>
            </div>
        </section>

        <section id="projects" class="projects-section py-5 bg-body-secondary text-light">
            <div class="container">
                <h2 class="text-center mb-5">Projects</h2>
                <div class="row g-4">
                    <?php foreach ($$Projects as $key => $Project): ?>
                        <article class="col-md-6 col-lg-4">
                            <div class="card project-card h-100 shadow-sm">
                                <div class="img-wrapper" style="background-image: url('<?= $Project->Image ?>')"></div>
                                <div class="card-body d-flex flex-column">
                                    <h3 class="card-title h5"><?= $Project->Title ?></h3>
                                    <p class="card-text flex-grow-1">
                                        <?= $Project->SmallDescription ?>
                                    </p>
                                    <a href="<?= $Project->Link ?>" class="btn btn-project mt-auto"
                                        target="_blank" rel="noopener">
                                        View on GitHub
                                    </a>
                                </div>
                            </div>
                        </article>
                    <?php endforeach ?>
                </div>
            </div>
        </section>

        <section id="contact-me" class="contact-section py-5 bg-body text-light">
            <div class="container">
                <h2 class="text-center mb-5">Contact Me</h2>

                <div class="row justify-content-center">
                    <!-- Contact Info -->
                    <div class="col-md-6 mb-4 d-flex align-items-center" style="background: transparent;">
                        <div class="w-100 text-start px-4">
                            <h4>Get in Touch</h4>
                            <p class="mb-1"><strong>Email:</strong> <a href="mailto:<?= $$ContactDetails->Email ?>"
                                    class="text-light"><?= $$ContactDetails->Email ?></a></p>
                            <p class="mb-1"><strong>Phone:</strong> <a href="tel:<?= $$ContactDetails->Phone ?>" class="text-light"><?= $$ContactDetails->Phone ?></a></p>
                            <p class="mb-0"><strong>Location:</strong> <?= $$ContactDetails->Location ?></p>
                            <hr class="border-light">
                            <!-- Social Media Links -->
                            <div class="d-flex">
                                <a href="<?= $$ContactDetails->Facebook ?>" class="text-light me-3 social-icon">
                                    <i class="bi bi-facebook fs-4"></i>
                                </a>
                                <a href="<?= $$ContactDetails->LinkedIn ?>" class="text-light me-3 social-icon">
                                    <i class="bi bi-linkedin fs-4"></i>
                                </a>
                                <a href="<?= $$ContactDetails->Instagram ?>" class="text-light me-3 social-icon">
                                    <i class="bi bi-instagram fs-4"></i>
                                </a>
                                <a href="<?= $$ContactDetails->GithubVrianta ?>" class="text-light me-3 social-icon">
                                    <i class="bi bi-github fs-4"></i>
                                </a>
                                <a href="<?= $$ContactDetails->Github ?>" class="text-light me-3 social-icon">
                                    <i class="bi bi-github fs-4"></i>
                                </a>
                            </div>
                        </div>
                    </div>

                    <!-- Contact Form -->
                    <div class="col-md-6 bg-body-secondary d-flex align-items-center">
                        <div class="p-4 rounded shadow-sm w-100">
                            <form id="contactForm" novalidate>
                                <div class="mb-3">
                                    <label for="name" class="form-label">Name <span class="text-danger">*</span></label>
                                    <input type="text" class="form-control" id="name" placeholder="Your name" required>
                                    <div class="invalid-feedback">Please enter your name.</div>
                                </div>

                                <div class="mb-3">
                                    <label for="email" class="form-label">Email <span
                                            class="text-danger">*</span></label>
                                    <input type="email" class="form-control" id="email"
                                        placeholder="your.email@example.com" required>
                                    <div class="invalid-feedback">Please enter a valid email address.</div>
                                </div>

                                <div class="mb-3">
                                    <label for="message" class="form-label">Message <span
                                            class="text-danger">*</span></label>
                                    <textarea class="form-control" id="message" rows="5"
                                        placeholder="Your message here..." required></textarea>
                                    <div class="invalid-feedback">Please enter your message.</div>
                                </div>

                                <button type="submit" class="btn btn-primary">Send Message</button>
                            </form>
                        </div>
                    </div>
                </div>
            </div>
        </section>

    </main>
    <!-- Footer -->
    <footer class="footer"> </footer>

    <!-- Bootstrap JS -->
    <script src="/js/Bootstrap/bootstrap.bundle.min.js" defer></script>
    <script src="/js/main.js" defer></script>

    <script>
        // Bootstrap 5 form validation
        (function() {
            'use strict'

            const form = document.getElementById('contactForm')

            form.addEventListener('submit', function(event) {
                if (!form.checkValidity()) {
                    event.preventDefault()
                    event.stopPropagation()
                } else {
                    // Here you can add your form submission code (e.g., AJAX)
                    event.preventDefault()
                    alert('Thank you for your message! I will get back to you soon.')
                    form.reset()
                }
                form.classList.add('was-validated')
            }, false)
        })()
    </script>
</body>

</html>

<script>
    const source = new EventSource("/hot-reload");

    source.onmessage = function(event) {
        if (event.data === "reload") {
            console.log("[LiveReload] Reloading page...");
            location.reload();
        }
    };

    source.onerror = function(err) {
        console.warn("[LiveReload] Disconnected from server", err);
    };
</script>