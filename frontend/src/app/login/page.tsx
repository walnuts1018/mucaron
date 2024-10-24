import Form from "next/form";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { MdLogin } from "react-icons/md";
import { useTranslation } from "react-i18next";

export default function Login() {
  const { t } = useTranslation("login");
  return (
    <div className="flex flex-col items-center justify-center h-screen">
      <h1>Login</h1>
      <Form action="/api/v1/login" className="flex flex-col">
        <label>
          {t("common.username", "ユーザー名")}
          <Input type="text" />
        </label>
        <label>
          {t("common.password", "パスワード")}
          <Input type="password" />
        </label>
        <Button type="submit">
          <MdLogin />
          {t("common.signin", "ログイン")}
        </Button>
      </Form>
    </div>
  );
}
