import { useEffect, useId, useRef, useState } from "react";

type MermaidDiagramProps = {
  source: string;
};

const mermaidBlockStart = /^(?:flowchart|graph|sequenceDiagram|classDiagram|stateDiagram(?:-v2)?|erDiagram|journey|gantt|pie|gitGraph|mindmap|timeline|quadrantChart|xychart-beta|block-beta|sankey-beta|architecture-beta)\b/;

type Mermaid = typeof import("mermaid").default;

let mermaidModule: Promise<Mermaid> | null = null;

export function isMermaidSource(source: string): boolean {
  return mermaidBlockStart.test(source.trim());
}

function loadMermaid(): Promise<Mermaid> {
  mermaidModule ??= import("mermaid").then(({ default: loadedMermaid }) => {
    loadedMermaid.initialize({
      startOnLoad: false,
      securityLevel: "strict",
      theme: "base",
      themeVariables: {
        background: "#101824",
        primaryColor: "#25352f",
        primaryTextColor: "#e2e5d8",
        primaryBorderColor: "#d9d9c4",
        lineColor: "#d9d9c4",
        secondaryColor: "#25352f",
        tertiaryColor: "#25352f",
        edgeLabelBackground: "#d9d9c4",
        fontFamily: "Cascadia Code, SFMono-Regular, Consolas, monospace",
        fontSize: "12px",
      },
    });
    return loadedMermaid;
  });
  return mermaidModule;
}

export function MermaidDiagram({ source }: MermaidDiagramProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const reactId = useId();
  const [error, setError] = useState<string | null>(null);
  const diagramId = `exercise-diagram-${reactId.replace(/[^a-zA-Z0-9]/g, "")}`;

  useEffect(() => {
    let cancelled = false;

    const renderDiagram = async () => {
      if (!containerRef.current) return;

      try {
        const mermaid = await loadMermaid();
        const { svg, bindFunctions } = await mermaid.render(
          diagramId,
          source.trim(),
        );
        if (cancelled || !containerRef.current) return;
        containerRef.current.innerHTML = svg;
        bindFunctions?.(containerRef.current);
        setError(null);
      } catch (renderError) {
        if (cancelled) return;
        setError(
          renderError instanceof Error
            ? renderError.message
            : "The diagram could not be rendered.",
        );
      }
    };

    void renderDiagram();
    return () => {
      cancelled = true;
    };
  }, [diagramId, source]);

  return (
    <div className="mermaid-diagram" role="img" aria-label="Exercise flowchart">
      <div className="mermaid-canvas" ref={containerRef} />
      {error ? (
        <details className="mermaid-error">
          <summary>Diagram unavailable</summary>
          <pre>{source}</pre>
          <small>{error}</small>
        </details>
      ) : null}
    </div>
  );
}
