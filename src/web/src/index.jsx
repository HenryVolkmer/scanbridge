import React, { useState, useEffect } from "react";
import { createRoot } from "react-dom/client";
import { CheckboxGroup, Checkbox, InlineLoading, Grid, Heading, Stack, Column, Form, Theme, TextInput, Button, InlineNotification } from "@carbon/react";
import "@carbon/styles/css/styles.css";

function ScanbridgeApp() {
  const [recipient, setRecipient] = useState("");
  const [colorMode, setColorMode] = useState(true);
  const [loading, setLoading] = useState(true);
  const [notification, setNotification] = useState({});

  useEffect(() => {
    async function fetchInitialValue() {
      try {
        const res = await fetch("/api/env");
        const data = await res.json();
        setRecipient(data.default_recipient ?? "");
      } catch (err) {
        setNotification({data: "Backend nicht erreichbar :(", kind: "error", title: "KO!"});
      } finally {
        setLoading(false);
      }
    }

    fetchInitialValue();
  }, []);

  const onSubmit = async (e) => {
    e.preventDefault();
    setNotification({});
    setLoading(true);
    try {
      const mode = colorMode === true ? "Color" : "Lineart";
      const res = await fetch(
        "/api/scan?mode=" + encodeURIComponent(mode) + "&recipient=" + encodeURIComponent(recipient)
      );
      const data = await res.json();
      if (!res.ok) {
        setNotification({data: data.Data, kind: "error", title: data.Title});
      } else {
        setNotification({data: data.Data, kind: "success", title: data.Title, url: data.url});
      }

    } catch (err) {
      console.error(err);
      setNotification({data: "Das hat nicht geklappt! :(", kind: "error", title: "KO!"});
    }
    setLoading(false);
  };

  return (
  <Theme theme="g10">
    <Grid>
      <Column sm={4}>
        <Form onSubmit={onSubmit}>
            <Stack gap={7}>
              <Heading>Scan starten</Heading>
              <p className="cds--body-long-01">
                Gib einen Empfänger für das PDF an und starte den Scanvorgang.
              </p>

              <CheckboxGroup
                legendText="Scanner Einstellungen"
              >
                <Checkbox
                  id="checkbox-color-enabled"
                  value={colorMode}
                  checked={colorMode}
                  labelText="Farbscan (langsamer)?"
                  onChange={(e) => setColorMode(e.target.checked)}
                />
              </CheckboxGroup>
              <TextInput
                id="simple-input"
                labelText="Empfänger E-Mail"
                placeholder="E-Mail eingeben ..."
                value={recipient}
                onChange={(e) => setRecipient(e.target.value)}
              />
              {notification.data && (
                <InlineNotification
                  kind={notification.kind}
                  statusIconDescription="notification"
                  subtitle={notification.data}
                  title={notification.title}
                  onClose={() => setNotification({})}
                  onCloseButtonClick={() => {}}
                />
              )}
              <Stack orientation="horizontal" gap={4}>
              {loading ? <InlineLoading
                status="active"
                description="scanne..."
              /> : <Button type="submit">bitti bitti Scani!</Button>}
              {notification?.url && <Button kind="secondary" onClick={() => window.location.href = notification.url}>Download</Button>}
            </Stack>

            </Stack>          
          </Form>
      </Column>
    </Grid>
  </Theme>
  );
}

const container = document.getElementById("root");
const root = createRoot(container);
root.render(<ScanbridgeApp />);
