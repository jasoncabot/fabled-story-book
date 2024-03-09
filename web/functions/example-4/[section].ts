import { EventContext, Request } from "@cloudflare/workers-types";
import { Env } from "..";

export async function onRequest(context: EventContext<Env, string, Request>) {
  const sectionId = context.params.section as string;

  const validRegex = /^[a-z0-9\-_]{1,59}\.jabl$/;
  if (!validRegex.test(sectionId)) {
    return new Response("Invalid section id", { status: 400 });
  }

  const source = "example-4";
  let jabl = await context.env.SECTIONS.get(`${source}:${sectionId}`);

  if (!jabl) {
    jabl = `{
      print("You've reached a point that doesn't exist!")
      print("This story needs a section with id '${sectionId}'.")
      choice("Start Again", {goto("entrypoint.jabl")})
      choice("Change Story", { set("system:source", 0) goto("entrypoint.jabl")})
    }`;
  }

  return new Response(jabl);
}
