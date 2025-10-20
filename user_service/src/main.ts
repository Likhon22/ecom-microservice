import express, { type Request, type Response } from "express";

const app = express();

const port = 4000;

app.get("/", (req: Request, res: Response) => {
  res.send("user service is running");
});

app.listen(port, () => {
  console.log("Server is running on port", port);
});
