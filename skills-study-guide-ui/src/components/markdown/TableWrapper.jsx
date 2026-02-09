export default function TableWrapper({ children }) {
  return (
    <div className="table-wrapper not-prose">
      <table className="w-full text-sm">
        {children}
      </table>
    </div>
  )
}
