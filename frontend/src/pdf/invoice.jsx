import React from 'react';
import { PDFViewer, Document, Page, Text } from '@react-pdf/renderer';
function MyDocument() {
  return (
    <Document>
      <Page size="A4">
        <Text>This is a sample PDF document.</Text>
      </Page>
    </Document>
  );
}

export default MyDocument;
