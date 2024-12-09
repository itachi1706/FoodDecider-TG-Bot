import mysql from 'mysql2/promise';

const dbConfig = {
  host: process.env.MYSQL_HOST || 'localhost',
  user: process.env.MYSQL_USER || 'root',
  port: Number(process.env.MYSQL_PORT) || 3306,
  password: process.env.MYSQL_PASSWORD || 'root',
  database: process.env.MYSQL_DATABASE || '',
  connectionLimit: Number(process.env.MYSQL_CONNECTION_LIMIT) || 5
};

console.log("Connecting to database with config: ", dbConfig.host, dbConfig.database);


const dbPool = mysql.createPool({
  connectionLimit: dbConfig.connectionLimit,
  host: dbConfig.host,
  user: dbConfig.user,
  port: dbConfig.port,
  password: dbConfig.password,
  database: dbConfig.database
});

export default dbPool;