import { z } from "zod";

const Login_FormSchema = z.object({
  email: z.email("Invalid email address"),
  password: z.string().min(1, "Required"),
  server: z.string().min(1, "Required"),
});

type Login_FormType = z.infer<typeof Login_FormSchema>;

export { Login_FormSchema, type Login_FormType };
