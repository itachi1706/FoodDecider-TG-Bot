import {TelegramUserData} from "@telegram-auth/server";

export type TelegramAuthenticationData = TelegramUserData & {
  is_admin: boolean;
  data_map: string;
}