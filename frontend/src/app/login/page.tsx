import Form from "next/form";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { MdLogin } from "react-icons/md";
import { useTranslations } from "next-intl";

export default function Login() {
  const t = useTranslations("Login");
  return (
    <div className="flex flex-col items-center justify-center h-screen">
      <h1>Login</h1>
      <Form action="/api/v1/login" className="flex flex-col">
        <label>
          {t("username")}
          <Input type="text" />
        </label>
        <label>
          {t("password")}
          <Input type="password" />
        </label>
        <Button type="submit">
          <MdLogin />
          ログイン
        </Button>
      </Form>
    </div>
  );
}
