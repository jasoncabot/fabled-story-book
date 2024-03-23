import React, { useEffect, useRef, useState } from "react";

const ConsoleText: React.FC<{ text: string }> = ({ text }) => {
  const consoleTextRef = useRef<HTMLParagraphElement>(null);
  const frame = useRef<number | null>(null);

  const [fullyRendered, setFullyRendered] = useState(false);

  const addCharacter = (char: string, renderStack: HTMLElement[]) => {
    let current = renderStack[renderStack.length - 1];

    // if the current node is the root paragraph then write this into a span
    if (current.nodeName === "P") {
      const span = document.createElement("div");
      current.appendChild(span);
      renderStack.push(span);
      current = span;
    }
    current.innerHTML += char;
  };

  const animateCharacter = (i: number, text: string, renderStack: HTMLElement[], animated: boolean) => {
    if (i >= text.length) {
      setFullyRendered(true);
      return;
    }

    const render = () => {
      const currentNode = renderStack[renderStack.length - 1];
      let char = text.charAt(i);
      switch (char) {
        case "/":
          if (currentNode.style.fontStyle === "italic") {
            renderStack.pop();
          } else {
            const italic = document.createElement("span");
            italic.style.fontStyle = "italic";
            currentNode.appendChild(italic);
            renderStack.push(italic);
          }
          break;
        case "*":
          if (currentNode.style.fontWeight === "bold") {
            renderStack.pop();
          } else {
            const bold = document.createElement("span");
            bold.style.fontWeight = "bold";
            currentNode.appendChild(bold);
            renderStack.push(bold);
          }
          break;
        case "_":
          if (currentNode.style.textDecoration === "underline") {
            renderStack.pop();
          } else {
            const underline = document.createElement("span");
            underline.style.textDecoration = "underline";
            currentNode.appendChild(underline);
            renderStack.push(underline);
          }
          break;
        case "`":
          if (currentNode.nodeName === "CODE") {
            renderStack.pop();
          } else {
            const code = document.createElement("code");
            code.classList.add("bg-gray-800", "text-harlequin-500", "p-1");
            currentNode.appendChild(code);
            renderStack.push(code);
          }
          break;
        case "\n":
          // if the currentNode is a header, pop it off the stack
          if (currentNode.nodeName.startsWith("H") || currentNode.nodeName === "DIV") {
            renderStack.pop();
          } else {
            const br = document.createElement("br");
            currentNode.appendChild(br);
          }
          break;
        case "#":
          let startedHeadersAt = i;
          while (text.charAt(i) === "#") {
            i++;
          }
          const size = i - startedHeadersAt;
          const header = document.createElement(`h${size}`);
          const sizes = ["text-4xl", "text-3xl", "text-2xl", "text-xl", "text-lg", "text-base"];
          header.classList.add(sizes[size - 1], "font-bold");
          currentNode.appendChild(header);
          renderStack.push(header);
          break;
        case "!":
          if (text.charAt(i + 1) === "[") {
            let imageText = "";
            i++;
            while (text.charAt(i) !== "(") {
              imageText += text.charAt(i);
              i++;
            }
            i++;
            let imageLink = "";
            while (text.charAt(i) !== ")") {
              imageLink += text.charAt(i);
              i++;
            }
            i++;
            imageText = imageText.substring(1, imageText.length - 1); // strip [ and ]
            const img = document.createElement("img");
            img.src = imageLink;
            img.title = imageText;
            img.alt = imageText;
            img.classList.add("w-full", "md:w-1/2", "m-auto", "h-auto");
            currentNode.appendChild(img);
          } else {
            addCharacter(char, renderStack);
          }
          break;
        default:
          addCharacter(char, renderStack);
          break;
      }

      animateCharacter(i + 1, text, renderStack, animated);
    };

    if (animated) {
      frame.current = requestAnimationFrame(render);
    } else {
      render();
    }
  };

  useEffect(() => {
    if (text && consoleTextRef.current) {
      consoleTextRef.current.innerHTML = "";
      animateCharacter(0, text, [consoleTextRef.current], true);
      return () => {
        if (frame.current) cancelAnimationFrame(frame.current);
      };
    }
  }, [text]);

  // can tap to skip the text animation
  const finishAnimation = () => {
    if (text && consoleTextRef.current && !fullyRendered) {
      consoleTextRef.current.innerHTML = "";
      if (frame.current) cancelAnimationFrame(frame.current);
      animateCharacter(0, text, [consoleTextRef.current], false);
    }
  };

  return (
    <div onClick={finishAnimation} className="flex flex-grow overflow-y-scroll overscroll-contain p-4">
      <p ref={consoleTextRef} onClick={finishAnimation} className="w-full text-justify font-mono text-harlequin-700"></p>
    </div>
  );
};

export default ConsoleText;
