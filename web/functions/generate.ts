import { Ai } from "@cloudflare/ai";
import { Request } from "@cloudflare/workers-types";
import { Env } from ".";

export async function onRequest({ request, env }: { request: Request; env: Env }) {
  // read the bearer token from the request
  const authHeader = request.headers.get("Authorization") ?? "";
  const [_, token] = authHeader.split(" ");
  if (token !== env.SECRET_KEY) {
    return new Response("Unauthorized", { status: 401 });
  }

  if (env.AI === undefined) {
    return new Response(
      `{
      print("AI not configured")
      choice("Try again", { goto("generate") })
    }`,
      { status: 200 },
    );
  }

  const ai = new Ai(env.AI);
  const input = {
    prompt:
      "In a realm shrouded in ancient lore and dark magic, a lone adventurer finds themselves ensnared within the labyrinthine depths of a mysterious dungeon. Each corridor and chamber is fraught with peril, housing puzzles of arcane origin and traps designed to thwart any attempts at escape. At the heart of the dungeon lies a promise of freedom, whispered by the enigmatic guardian who watches from the shadows. Yet, as the adventurer delves deeper, they realize that dark forces conspire to ensure that the secrets of the dungeon remain forever concealed. With their courage tested and their wits challenged, the adventurer must navigate the treacherous maze, unravel its mysteries, and confront the guardians that stand in the way of their ultimate escape. Craft a gripping hook that plunges the player into the thrilling journey of a solitary hero against the odds.",
  };
  const answer = (await ai.run("@cf/meta/llama-2-7b-chat-int8", input)).response;
  const jabl = `{
     ${answer
       .split("\n")
       .map((line) => `print("${line.replace(/^"|"$/g, "")}")`)
       .join("\n")}
     choice("Continue", { goto("generate") })
  }`;
  return new Response(jabl);
}
