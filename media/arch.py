from diagrams import Cluster, Diagram, Edge
from diagrams.aws.compute import EC2
from diagrams.aws.database import RDS
from diagrams.aws.network import ELB, VPC, NATGateway, InternetGateway
from diagrams.aws.management import SystemsManager as SSM
from diagrams.onprem.client import Users
from diagrams.onprem.network import Internet

with Diagram("3 Tier Web Module", show=False):
    users = Users("Users")
    internet = Internet("Internet")
    with Cluster("Region Sa-East-1"):
        ssm = SSM("SSM Management")
        with Cluster("VPC"):
            with Cluster("subnet/public"):
                igw = InternetGateway("InternetGateway")
                lb = ELB("lb")

            with Cluster("subnet/private"):
                with Cluster("App"):
                    auto_scaling_group = [EC2("app1"), EC2("app2")]
                with Cluster("Database"):
                    database = RDS("app_db")
                natgw = NATGateway("NATGW")

    users >> internet >> igw >> lb >> natgw
    natgw >> auto_scaling_group
    auto_scaling_group - Edge(style="dotted") - database
    ssm >> Cluster("subnet/private")