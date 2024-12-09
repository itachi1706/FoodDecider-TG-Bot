import db from "@/utils/database";
import {RowDataPacket} from "mysql2";

export async function checkIsAdmin(id : number) {
  const stmt = 'SELECT * FROM admins WHERE telegram_id = ? AND status = "A"';

  const [rows] = await db.execute<RowDataPacket[]>(stmt, [id]);

  if (rows.length == 0) {
    console.log("User not found in database. Not an admin");
    return false
  }

  console.log("User found in database. Is an admin");
  return true;
}