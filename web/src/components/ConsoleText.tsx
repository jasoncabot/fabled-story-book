import React, { useEffect, useRef, useState } from "react";

const ConsoleText: React.FC<{ text: string }> = ({ text }) => {
  const consoleTextRef = useRef<HTMLParagraphElement>(null);
  const frame = useRef<number | null>(null);

  const [fullyRendered, setFullyRendered] = useState(false);

  const addElement = (element: HTMLElement, renderStack: HTMLElement[]) => {
    let currentNode = renderStack[renderStack.length - 1];

    if (currentNode.nodeName === "P") {
      const div = document.createElement("div");
      div.classList.add("mb-4");

      currentNode.appendChild(div);
      renderStack.push(div);

      currentNode = div;
    }

    currentNode.appendChild(element);
    renderStack.push(element);
  };

  const addCharacter = (char: string, renderStack: HTMLElement[]) => {
    let current = renderStack[renderStack.length - 1];

    // if the current node is the root paragraph then write this into a span
    if (current.nodeName === "P") {
      const div = document.createElement("div");
      div.classList.add("mb-4");
      current.appendChild(div);
      renderStack.push(div);
      current = div;
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
            addElement(italic, renderStack);
          }
          break;
        case "*":
          if (currentNode.style.fontWeight === "bold") {
            renderStack.pop();
          } else {
            const bold = document.createElement("span");
            bold.style.fontWeight = "bold";
            addElement(bold, renderStack);
          }
          break;
        case "_":
          if (currentNode.style.textDecoration === "underline") {
            renderStack.pop();
          } else {
            const underline = document.createElement("span");
            underline.style.textDecoration = "underline";
            addElement(underline, renderStack);
          }
          break;
        case "`":
          if (currentNode.nodeName === "CODE") {
            renderStack.pop();
          } else {
            const code = document.createElement("code");
            code.classList.add("bg-gray-800", "text-harlequin-500", "p-1");
            addElement(code, renderStack);
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
        case "|":
          let tableText = "";
          while (text.charAt(i) == "|") {
            let j = i;
            while (text.charAt(j) !== "\n") {
              tableText += text.charAt(j);
              j++;
            }
            if (text.charAt(j + 1) == "|") {
              j++;
            }
            tableText += "\n";
            i = j;
          }
          tableText = tableText.substring(0, tableText.length - 1);

          const rows = tableText.split("\n");
          const table = document.createElement("table");
          table.classList.add("w-full", "mb-4", "border-collapse");
          for (let row of rows) {
            const tr = document.createElement("tr");
            const cells = row.split("|").map((cell) => cell.trim());
            cells.shift();
            cells.pop();
            for (let cell of cells) {
              const td = document.createElement("td");
              td.classList.add("border", "border-harlequin-500", "px-4", "py-2");
              td.textContent = cell;
              tr.appendChild(td);
            }
            table.appendChild(tr);
          }
          currentNode.appendChild(table);
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
            img.classList.add("w-full", "md:w-1/2", "m-auto", "h-auto", "rounded-lg", "shadow-lg", "mb-4");
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
      setFullyRendered(false);
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
      <p ref={consoleTextRef} className="w-full font-mono text-harlequin-700"></p>
    </div>
  );
};

export default ConsoleText;
