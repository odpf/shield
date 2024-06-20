import { V1Beta1User } from "@raystack/frontier";
import type { ColumnDef } from "@tanstack/react-table";
import { Link } from "react-router-dom";

export const getColumns: () => ColumnDef<V1Beta1User, any>[] = () => {
  return [
    {
      accessorKey: "id",
      header: "ID",
      filterVariant: "text",
      cell: ({ row, getValue }) => {
        return <Link to={`/users/${row.getValue("id")}`}>{getValue()}</Link>;
      },
    },
    {
      header: "Name",
      accessorKey: "name",
      filterVariant: "text",
      cell: (info) => info.getValue(),
    },
    {
      header: "Email",
      accessorKey: "email",
      filterVariant: "text",
      cell: (info) => info.getValue(),
      footer: (props) => props.column.id,
    },
    {
      header: "Created At",
      accessorKey: "created_at",
      meta: {
        headerFilter: false,
      },
      cell: (info) =>
        new Date(info.getValue() as Date).toLocaleString("en", {
          month: "long",
          day: "numeric",
          year: "numeric",
        }),
      filterVariant: "date",
      footer: (props) => props.column.id,
    },
  ];
};
