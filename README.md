This project will follow this architecture design:

![](media/architecture.png)

With this in mind, this project has the intention of doing the following:

- [ ] Creating a new VPC with a private and a public subnet
- [ ] Placing an RDS MySQL inside the private Subnet
- [ ] Creating an AutoScaling Group in the public subnet with a Load Balancer
- [ ] Each instance will serve an static file containing the version of the Database
- [ ] Write an Unit Test to check if the version returned on the static file is the correct one
- [ ] Set all this project using Github Actions