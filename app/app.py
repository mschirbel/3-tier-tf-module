from flask import Flask
import mysql.connector
from mysql.connector import errorcode
import os

app = Flask(__name__)

@app.route('/')
def hello():
    rds_main = {
        'user': os.environ['RDS_USER'],
        'password': os.environ['RDS_PWRD'],
        'host': os.environ['RDS_HOST'],
        'database': os.environ['RDS_BASE'],
        'raise_on_warnings': True
    }
    try:
        conn = mysql.connector.connect(**rds_main)
        cursor = conn.cursor()  
        db_version = cursor.execute("SELECT VERSION()")
        return db_version
        conn.close()
    except mysql.connector.Error as err:
        print(err)
        raise ValueError("Some Error in DB Connection")
    