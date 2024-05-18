- First Lets find the API

- Then lets create env file with the needed secrets/variables

- Now lets learn a bit about RabbitMQ

- Create the docker-compose file

- Go through the RabbitMQ get started guide

- Update .env file with the data needed for accessing rabbitmq

- Now lets structure the program and write the boilerplate messaging system

- Write the logic for the command line app to make it take params, query and return
    - each subscriber reads the task, performs it, saves the data and responds to the client (RPC request/reply pattern)

- Create the business logic code for getting the daily movie suggestions
    - get a genre from the user as a flag
    - or get random genres

- command line app is the publisher of the tasks

- through the cmd you can say how many subscribers you want

- Write the code to convert images to ANSI RGB control codes and Unicode block graphics characters for the terminal

- Integrate it into the application

- Create the dockerfile for dockerizing this application and integrate it into the docker-compose file