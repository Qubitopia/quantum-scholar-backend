import 'bootstrap/dist/css/bootstrap.min.css';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Classroom from './pages/classroom.jsx';
import Home from './pages/home.jsx';
import Settings from './pages/settings.jsx';

function App() {
  // Example: userType could be determined from authentication or context
  const userType = "teacher"; // or "teacher"

  return (
    <Router>
      <Routes>
        {userType === "admin" ? (
          <>
            <Route path="/" element={<Home />} />
            <Route path="/classroom" element={<Classroom />} />
            <Route path="/settings" element={<Settings />} />
          </>
        ) : (
          <>
            <Route path="/" element={<Home />} />
            <Route path="/classroom" element={<Classroom />} />
          </>
        )}
      </Routes>
    </Router>
  );
}

export default App;
