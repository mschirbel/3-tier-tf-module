from flask import Flask
from flask import jsonify
import mysql.connector
from mysql.connector import errorcode
import os
from ec2_metadata import ec2_metadata
from dotenv import load_dotenv

app = Flask(__name__)
load_dotenv()

@app.route('/')
def hello():
    rds_main = {
        'user': os.getenv("RDS_USER"),
        'password': os.getenv("RDS_PWRD"),
        'host':os.getenv("RDS_HOST"),
        'database': os.getenv("RDS_BASE"),
        'raise_on_warnings': True
    }
    try:
        conn = mysql.connector.connect(**rds_main)
        cursor = conn.cursor()  
        cursor.execute("SELECT VERSION()")
        db_version = cursor.fetchone()
        return jsonify(
            database_version=db_version,
            region=ec2_metadata.region,
            unique_id=os.getenv("UNIQUE_ID ")
        )
        conn.close()
    except mysql.connector.Error as err:
        print(err)
        raise ValueError("Some Error in DB Connection")
    except TypeError as e:
        print(e)
        exit(1)
    