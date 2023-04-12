export async function update(
  url: string,
  { arg }: { arg: Record<string, string> }
) {
  await fetch(url, {
    method: "POST",
    headers: {},
    body: JSON.stringify(arg),
  });
}

export async function updateOrganisation(
  url: string,
  { arg }: { arg: Record<string, string> }
) {
  await fetch(url, {
    method: "POST",
    headers: {},
    body: JSON.stringify(arg),
  });
}
