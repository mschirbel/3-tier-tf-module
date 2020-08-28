This project will follow this architecture design:

![](media/architecture.png)

With this in mind, this project has the intention of doing the following:

- [x] Creating a new VPC with a private and a public subnet
- [x] Placing an RDS MySQL inside the private Subnet
- [x] The RDS must have a random password
- [x] The RDS password will be stored in a SSM Parameter Store
- [x] Creating an AutoScaling Group in the public subnet with a Load Balancer
- [x] Each instance will serve an api which returns the database version
- [ ] Write an Unit Test to check if the version returned on the api is the correct one
- [ ] Set all this project using Github Actions