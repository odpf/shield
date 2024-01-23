import { Flex, Grid, Text } from "@raystack/apsara";
import { useFrontier } from "@raystack/frontier/react";
import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import PageHeader from "~/components/page-header";
import { User } from "~/types/user";

export default function UserDetails() {
  const { client } = useFrontier();
  let { userId } = useParams();
  const [user, setUser] = useState<User>();

  useEffect(() => {
    async function getProject() {
      const {
        // @ts-ignore
        data: { user },
      } = await client?.frontierServiceGetUser(userId ?? "");
      setUser(user);
    }
    getProject();
  }, [userId]);

  const pageHeader = {
    title: "Users",
    breadcrumb: [
      {
        href: `/users`,
        name: `Users list`,
      },
      {
        href: `/users/${user?.id}`,
        name: `${user?.email}`,
      },
    ],
  };

  return (
    <Flex
      direction="column"
      gap="large"
      style={{
        width: "320px",
        height: "calc(100vh - 60px)",
        borderLeft: "1px solid var(--border-base)",
      }}
    >
      <PageHeader title={pageHeader.title} breadcrumb={pageHeader.breadcrumb} />
      <Flex direction="column" gap="large" style={{ padding: "0 24px" }}>
        <Grid columns={2} gap="small">
          <Text size={1}>Email</Text>
          <Text size={1}>{user?.email}</Text>
        </Grid>
        <Grid columns={2} gap="small">
          <Text size={1}>Created At</Text>
          <Text size={1}>
            {new Date(user?.created_at as Date).toLocaleString("en", {
              month: "long",
              day: "numeric",
              year: "numeric",
            })}
          </Text>
        </Grid>
      </Flex>
    </Flex>
  );
}