import type { CompareContent } from '../../content/landingContent';

export default function CompareTable({ content }: { content: CompareContent }) {
  return (
    <section className="py-16 bg-gray-50">
      <div className="max-w-6xl mx-auto px-4">
        <h2 className="text-2xl md:text-3xl font-bold text-center mb-8">{content.title}</h2>
        <div className="overflow-x-auto">
          <table className="min-w-full border border-gray-200 rounded-lg bg-white">
            <thead>
              <tr>
                {content.columns.map((col: string, i: number) => (
                  <th key={i} className="px-4 py-2 text-left font-semibold text-gray-700 border-b border-gray-200 bg-gray-100">
                    {col}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {content.rows.map((row: string[], i: number) => (
                <tr key={i} className="hover:bg-gray-50">
                  {row.map((cell: string, j: number) => (
                    <td key={j} className="px-4 py-2 text-sm text-gray-800 border-b border-gray-100">
                      {cell}
                    </td>
                  ))}
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </section>
  );
} 