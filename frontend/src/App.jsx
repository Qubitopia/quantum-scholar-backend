import 'bootstrap/dist/css/bootstrap.min.css';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Classroom from './pages/classroom.jsx';
import Home from './pages/home.jsx';

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/classroom" element={<Classroom />} />
      </Routes>
    </Router>
  );
}

export default App;
