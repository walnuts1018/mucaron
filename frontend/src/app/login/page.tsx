export default function Login() {
  return (
    <div>
      <h1>Login</h1>
      <form>
        <label>
          ユーザー名
          <input type="text" />
        </label>
        <label>
          Password
          <input type="password" />
        </label>
        <button type="submit">Login</button>
      </form>
    </div>
  );
}
