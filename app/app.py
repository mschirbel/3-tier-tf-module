from flask import Flask
import mysql.connector
from mysql.connector import errorcode

app = Flask(__name__)

@app.route('/')
def hello():
    rds_main = {
        'user': os.environ['RDS_USER'],
        'pwrd': os.environ['RDS_PWRD'],
        'host': os.environ['RDS_HOST'],
        'base': os.environ['RDS_BASE'],
        'raise_on_warnings': True
    }
    try:
        conn = mysql.connector(**rds_main)
        return 'Connected with RDS!'
    except mysql.connector.Error as err:
        return err
    