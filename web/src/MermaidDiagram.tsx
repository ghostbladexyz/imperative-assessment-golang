import { useEffect, useId, useRef, useState } from "react";

type MermaidDiagramProps = {
  source: string;
};

const mermaidFlowchartStart = /^(?:flowchart|graph)\s+(?:TB|TD|BT|RL|LR)(?:\s|$)/;

type Mermaid = typeof import("mermaid").default;

let mermaidModule: Promise<Mermaid> | null = null;
let mermaidRenderSequence = 0;

export function isMermaidSource(source: string): boolean {
  return mermaidFlowchartStart.test(source.trim());
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
      const container = containerRef.current;
      if (!container) return;

      container.replaceChildren();
      setError(null);

      try {
        const mermaid = await loadMermaid();
        const renderId = `${diagramId}-${++mermaidRenderSequence}`;
        const { svg, bindFunctions } = await mermaid.render(
          renderId,
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
