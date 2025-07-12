package main

/*
This is a handler of my overall package
This will do following things

1. Create Default Application
2. Create Default Controller
3. Create Default Module
4. Controll the migration
5. Loop Over all the files to check for update and if any file updated then it will restart the application auto matically

----- Command line arguments
1. --new-application app_name / -na app_name (if will make sure the current directory is not a application folder by checking coltroller/model/views folders)
2. --new-controller controller_name / -nc (if controller for the same is not present then it will create it)
3. --new-model model_name / -nm (create new model if that is not already create - only the ID attribute will be there in the default mpdel)
4. --migrate-model / -mm (check the current model if the build is true in config then the it will confirm before doing any changes in the model)
5. --migrate-component / -mc (Migrate Components with the DB and ask for deletion or modification in Build mode true)
6. --start-server / -ss (start the server with the configaration mentioned in the config folder)

-------------------------------------------------------------------------------------------------------------------------------------------------

This file will be usefull for future expansion of extension and UI based Modifications
TODO:
IDK :P
*/
