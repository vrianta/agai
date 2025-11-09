<header class="sticky-top">
    <nav class="navbar navbar-expand-lg navbar-light bg-body">
        <div class="container">
            <a class="navbar-brand" href=""><?= $$Heading ?></a>
            <button class="navbar-toggler" type="button" data-bs-toggle="collapse"
                data-bs-target="#navbarSupportedContent" aria-controls="navbarSupportedContent"
                aria-expanded="false" aria-label="Toggle navigation">
                <span class="navbar-toggler-icon"></span>
            </button>
            <!-- Navigation Items which are going to be shown in the navigation bar -->
            <div class="collapse navbar-collapse" id="navbarSupportedContent">
                <ul class="navbar-nav ms-auto mb-2 mb-lg-0">
                    <?php foreach ($$NavItems as $key => $NavItem): ?>
                        <?php if ($NavItem->DropDown): ?>
                            <li class="nav-item dropdown">
                                <a class="nav-link dropdown-toggle" href="<?= $NavItem->Href ?>" role="button"
                                    data-bs-toggle="dropdown" aria-expanded="false"><?= $NavItem->Name ?></a>
                                <ul class="dropdown-menu">
                                    <?php foreach ($NavItem->DropDown as $key => $DropDownItem): ?>
                                        <li><a class="dropdown-item"
                                                href="<?= $DropDownItem->Href ?>"><?= $DropDownItem->Name ?></a></li>
                                    <?php endforeach ?>
                                </ul>
                            </li>
                        <?php elseif ($NavItem->Disabled): ?>
                            <li class="nav-item">
                                <a class="nav-link disabled" aria-disabled="true"><?= $NavItem->Name ?></a>
                            </li>
                        <?php else: ?>
                            <li class="nav-item">
                                <a class="nav-link active" aria-current="page"
                                    href="<?= $NavItem->Href ?>"><?= $NavItem->Name ?></a>
                            </li>
                        <?php endif ?>
                    <?php endforeach ?>
                </ul>
            </div>
        </div>
    </nav>

</header>