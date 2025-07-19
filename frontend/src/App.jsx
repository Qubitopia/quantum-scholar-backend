import { PDFViewer } from '@react-pdf/renderer';
import MyDocument from './pdf/invoice.jsx';
import './App.css'


function App() {
  return (
    <div className='app-container'>
      <div className='left-section'>
        <h1>PDF Viewer</h1>
        <p>View the PDF document above.</p>
      </div>
      <div className='right-section'>
        <PDFViewer width="100%" height={500}>
          <MyDocument />
        </PDFViewer>
      </div>
    </div>


  );
}

export default App;
