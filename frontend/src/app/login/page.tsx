import Form from "next/form";

export default function Login() {
  return (
    <div className="flex flex-col items-center justify-center h-screen">
      <h1>Login</h1>
      <Form action="/api/v1/login">
        <label>
          ユーザー名
          <input type="text" />
        </label>
        <label>
          Password
          <input type="password" />
        </label>
        <button type="submit">Login</button>
      </Form>
    </div>
  );
}
