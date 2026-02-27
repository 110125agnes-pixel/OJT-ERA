import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import "./EmployeeProfiling.css";
import { itemService } from "./services/api";
import ItemForm from "./components/ItemForm";
import ItemEditForm from "./components/ItemEditForm";

function EmployeeProfiling() {
  const navigate = useNavigate();
  const [employees, setEmployees] = useState([]);
  const [newEmployee, setNewEmployee] = useState({
    caseNo: "",
    hospitalNo: "",
    lastname: "",
    firstname: "",
    middlename: "",
    suffix: "",
    birthdate: "",
    age: "",
    room: "",
    admissionDate: "",
    dischargeDate: "",
    sex: "",
    height: "",
    weight: "",
    complaint: "",
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [searchTerm, setSearchTerm] = useState("");
  const [editItem, setEditItem] = useState(null);
  const [editLoading, setEditLoading] = useState(false);

  const sexOptions = ["", "Male", "Female", "Other"];
  const civilStatusOptions = ["", "Single", "Married", "Divorced", "Widowed"];

  // Input sanitizers
  const onlyDigits = (val) => (val || "").toString().replace(/\D+/g, "");
  const onlyLetters = (val) =>
    (val || "").toString().replace(/[^a-zA-Z\s]+/g, "");

  // Generate a random alphanumeric string (uppercase)
  const generateRandomCaseNo = (len = 8) => {
    const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789";
    let out = "";
    for (let i = 0; i < len; i++)
      out += chars.charAt(Math.floor(Math.random() * chars.length));
    return out;
  };

  // Generate a unique caseNo not present in the provided list
  const generateUniqueCaseNo = (list, len = 8) => {
    const existing = new Set(
      (list || [])
        .map((e) => (e.caseNo || "").toString().trim())
        .filter(Boolean),
    );
    let attempt = 0;
    while (attempt < 1000) {
      const candidate = generateRandomCaseNo(len);
      if (!existing.has(candidate)) return candidate;
      attempt += 1;
    }
    // fallback to timestamp-based unique value
    return `CN${Date.now()}`;
  };

  // Date helpers: format values for backend MySQL fields
  // - `toSQLDate` converts a date-like string to `YYYY-MM-DD` which
  //   matches MySQL `DATE` column expectations.
  // - `toSQLDateTime` converts to `YYYY-MM-DD HH:MM:SS` which matches
  //   MySQL `DATETIME` columns. These helpers prevent invalid/ambiguous
  //   date formats being sent from the browser (which can trigger 500s).
  const toSQLDate = (dateStr) => {
    if (!dateStr) return "";
    const d = new Date(dateStr);
    if (isNaN(d)) return dateStr;
    const yyyy = d.getFullYear();
    const mm = String(d.getMonth() + 1).padStart(2, "0");
    const dd = String(d.getDate()).padStart(2, "0");
    return `${yyyy}-${mm}-${dd}`;
  };

  const toSQLDateTime = (dateStr) => {
    if (!dateStr) return "";
    const d = new Date(dateStr);
    if (isNaN(d)) return dateStr;
    const yyyy = d.getFullYear();
    const mm = String(d.getMonth() + 1).padStart(2, "0");
    const dd = String(d.getDate()).padStart(2, "0");
    const hh = String(d.getHours()).padStart(2, "0");
    const min = String(d.getMinutes()).padStart(2, "0");
    const ss = String(d.getSeconds()).padStart(2, "0");
    return `${yyyy}-${mm}-${dd} ${hh}:${min}:${ss}`;
  };

  useEffect(() => {
    fetchEmployees();
  }, []);

  const fetchEmployees = async () => {
    try {
      setLoading(true);
      const data = await itemService.getAllItems();
      // Filter out rows that are clearly not patient records (e.g. library rows)
      // Heuristic: require at least a firstname or lastname to be present.
      const filtered = (data || []).filter((d) => {
        const ln = (d.lastname || "").toString().trim();
        const fn = (d.firstname || "").toString().trim();
        return ln.length > 0 || fn.length > 0;
      });
      setEmployees(filtered);
      // set a generated unique caseNo for the form using the filtered list
      setNewEmployee((prev) => ({
        ...prev,
        caseNo: generateUniqueCaseNo(filtered),
      }));
      setError("");
    } catch (err) {
      setError(
        "Failed to fetch employees: " +
          (err.response?.data?.error || err.message),
      );
      console.error("Error fetching employees:", err);
    } finally {
      setLoading(false);
    }
  };

  const addEmployee = async (e) => {
    e.preventDefault();
    if (!newEmployee.lastname.trim() || !newEmployee.firstname.trim()) {
      setError("Lastname and Firstname are required");
      return;
    }
    try {
      setLoading(true);
      // Prepare payload: explicitly normalize date/time fields so the
      // backend receives consistent formats (DATE and DATETIME). This
      // avoids server-side parsing errors when the user edits twice.
      const payload = {
        ...newEmployee,
        // MySQL DATE
        birthdate: toSQLDate(newEmployee.birthdate),
        // MySQL DATETIME (send null if empty)
        admissionDate: newEmployee.admissionDate
          ? toSQLDateTime(newEmployee.admissionDate)
          : null,
        dischargeDate: newEmployee.dischargeDate
          ? toSQLDateTime(newEmployee.dischargeDate)
          : null,
      };
      const data = await itemService.createItem(payload);
      const updated = [...employees, data];
      setEmployees(updated);
      setNewEmployee({
        hospitalNo: "",
        lastname: "",
        firstname: "",
        middlename: "",
        suffix: "",
        birthdate: "",
        age: "",
        room: "",
        admissionDate: "",
        dischargeDate: "",
        sex: "",
        height: "",
        weight: "",
        complaint: "",
      });
      setError("");
      // set next generated caseNo for the next entry
      setNewEmployee((prev) => ({
        ...prev,
        caseNo: generateUniqueCaseNo(updated),
      }));
    } catch (err) {
      setError(
        "Failed to add employee: " + (err.response?.data?.error || err.message),
      );
      console.error("Error adding employee:", err);
    } finally {
      setLoading(false);
    }
  };

  const deleteEmployee = async (id) => {
    const confirmDelete = window.confirm(
      "Are you sure you want to delete this employee?",
    );
    if (!confirmDelete) return;

    try {
      setLoading(true);
      await itemService.deleteItem(id);
      const remaining = employees.filter((emp) => emp.id !== id);
      setEmployees(remaining);
      // regenerate caseNo after deletion to avoid collisions
      setNewEmployee((prev) => ({
        ...prev,
        caseNo: generateUniqueCaseNo(remaining),
      }));
      setError("");
    } catch (err) {
      setError("Failed to delete employee: " + err.message);
      console.error("Error deleting employee:", err);
    } finally {
      setLoading(false);
    }
  };

  const openEdit = (emp) => {
    setEditItem({ ...emp });
  };

  const handleEditChange = (updated) => {
    setEditItem(updated);
  };

  const saveEdit = async () => {
    if (!editItem) return;
    try {
      setEditLoading(true);
      // Prepare payload for update: normalize dates so the PUT body
      // contains the same formats the backend expects. This reduces
      // chance of server errors when updating an already-edited record.
      const payload = {
        ...editItem,
        birthdate: toSQLDate(editItem.birthdate),
        admissionDate: toSQLDateTime(editItem.admissionDate),
        dischargeDate: toSQLDateTime(editItem.dischargeDate),
      };
      const updated = await itemService.updateItem(editItem.id, payload);
      setEmployees(employees.map((e) => (e.id === updated.id ? updated : e)));
      setEditItem(null);
      setError("");
    } catch (err) {
      setError(
        "Failed to update employee: " +
          (err.response?.data?.error || err.message),
      );
      console.error("Error updating employee:", err);
    } finally {
      setEditLoading(false);
    }
  };

  const cancelEdit = () => {
    setEditItem(null);
  };

  const viewPatient = (patient) => {
    navigate(`/patient/${patient.id}`);
  };

  const filteredEmployees = employees.filter(
    (emp) =>
      emp.firstname?.toLowerCase().includes(searchTerm.toLowerCase()) ||
      emp.lastname?.toLowerCase().includes(searchTerm.toLowerCase()),
  );

  return (
    <div className="profiling-container">
      <h2>👥 Patient Management</h2>
      <p className="subtitle">Create and manage patient profiles</p>

      {error && <div className="error">{error}</div>}

      <form onSubmit={addEmployee} className="add-form">
        <h3>Add New Patient</h3>
        <div className="form-grid">
          <div className="form-field">
            <label>Case No. *</label>
            <input
              name="caseNo"
              type="text"
              value={newEmployee.caseNo}
              readOnly
              required
            />
          </div>
          <div className="form-field">
            <label>Hospital No.</label>
            <input
              name="hospitalNo"
              type="text"
              value={newEmployee.hospitalNo}
              onChange={(e) =>
                setNewEmployee({
                  ...newEmployee,
                  hospitalNo: onlyDigits(e.target.value),
                })
              }
            />
          </div>
          <div className="form-field">
            <label>Last Name *</label>
            <input
              name="lastname"
              type="text"
              value={newEmployee.lastname}
              onChange={(e) =>
                setNewEmployee({
                  ...newEmployee,
                  lastname: onlyLetters(e.target.value),
                })
              }
              required
            />
          </div>
          <div className="form-field">
            <label>First Name *</label>
            <input
              name="firstname"
              type="text"
              value={newEmployee.firstname}
              onChange={(e) =>
                setNewEmployee({
                  ...newEmployee,
                  firstname: onlyLetters(e.target.value),
                })
              }
              required
            />
          </div>

          <div className="form-field">
            <label>Middle Name</label>
            <input
              name="middlename"
              type="text"
              value={newEmployee.middlename}
              onChange={(e) =>
                setNewEmployee({
                  ...newEmployee,
                  middlename: onlyLetters(e.target.value),
                })
              }
            />
          </div>
          <div className="form-field">
            <label>Suffix</label>
            <input
              name="suffix"
              type="text"
              value={newEmployee.suffix}
              onChange={(e) =>
                setNewEmployee({
                  ...newEmployee,
                  suffix: onlyLetters(e.target.value),
                })
              }
            />
          </div>
          <div className="form-field">
            <label>Birthdate *</label>
            <input
              name="birthdate"
              type="date"
              value={newEmployee.birthdate}
              onChange={(e) =>
                setNewEmployee({ ...newEmployee, birthdate: e.target.value })
              }
              required
            />
          </div>
          <div className="form-field">
            <label>Age</label>
            <input
              name="age"
              type="text"
              value={newEmployee.age}
              onChange={(e) =>
                setNewEmployee({
                  ...newEmployee,
                  age: onlyDigits(e.target.value),
                })
              }
            />
          </div>

          <div className="form-field">
            <label>Room</label>
            <input
              name="room"
              type="text"
              value={newEmployee.room}
              onChange={(e) =>
                setNewEmployee({
                  ...newEmployee,
                  room: onlyDigits(e.target.value),
                })
              }
            />
          </div>
          <div className="form-field">
            <label>Admission Date</label>
            <input
              name="admissionDate"
              type="datetime-local"
              value={newEmployee.admissionDate}
              onChange={(e) =>
                setNewEmployee({
                  ...newEmployee,
                  admissionDate: e.target.value,
                })
              }
            />
          </div>
          <div className="form-field">
            <label>Discharge Date</label>
            <input
              name="dischargeDate"
              type="datetime-local"
              value={newEmployee.dischargeDate}
              onChange={(e) =>
                setNewEmployee({
                  ...newEmployee,
                  dischargeDate: e.target.value,
                })
              }
            />
          </div>
          <div className="form-field">
            <label>Sex</label>
            <select
              name="sex"
              value={newEmployee.sex}
              onChange={(e) =>
                setNewEmployee({ ...newEmployee, sex: e.target.value })
              }
            >
              <option value="">Select Sex</option>
              <option value="Male">Male</option>
              <option value="Female">Female</option>
            </select>
          </div>

          <div className="form-field">
            <label>Height (cm)</label>
            <input
              name="height"
              type="text"
              value={newEmployee.height}
              onChange={(e) =>
                setNewEmployee({
                  ...newEmployee,
                  height: onlyDigits(e.target.value),
                })
              }
            />
          </div>
          <div className="form-field">
            <label>Weight (kg)</label>
            <input
              name="weight"
              type="text"
              value={newEmployee.weight}
              onChange={(e) =>
                setNewEmployee({
                  ...newEmployee,
                  weight: onlyDigits(e.target.value),
                })
              }
            />
          </div>
          <div className="form-field">
            <label>Complaint</label>
            <input
              name="complaint"
              type="text"
              value={newEmployee.complaint}
              onChange={(e) =>
                setNewEmployee({
                  ...newEmployee,
                  complaint: onlyLetters(e.target.value),
                })
              }
            />
          </div>
          <div className="form-field">
            <label>&nbsp;</label>
            <button type="submit" disabled={loading} className="submit-btn">
              {loading ? "Adding..." : "Add Patient"}
            </button>
          </div>
        </div>
      </form>

      <div className="search-section">
        <input
          type="text"
          placeholder="🔍 Search by name..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          className="search-input"
        />
      </div>
      {/* Edit Modal */}
      {editItem && (
        <div className="modal-overlay">
          <div className="modal-content">
            <div className="modal-header">
              <h3>Edit Patient</h3>
              <button className="close-btn" onClick={cancelEdit}>
                &times;
              </button>
            </div>
            <div className="patient-details">
              <ItemEditForm
                item={editItem}
                onChange={handleEditChange}
                onSave={saveEdit}
                onCancel={cancelEdit}
                loading={editLoading}
                sexOptions={sexOptions}
                civilStatusOptions={civilStatusOptions}
              />
            </div>
          </div>
        </div>
      )}

      <div className="employees-section">
        <h3>Patient List ({filteredEmployees.length})</h3>
        {loading && <p>Loading...</p>}
        <div className="employees-table-container">
          <table className="employees-table">
            <thead>
              <tr>
                <th>Case No.</th>
                <th>Full Name</th>
                <th>Birthdate</th>
                <th>Room</th>
                <th>Admission Date</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {filteredEmployees.map((emp) => (
                <tr key={emp.id}>
                  <td>{emp.caseNo}</td>
                  <td>
                    {`${emp.lastname}, ${emp.firstname} ${emp.middlename || ""} ${emp.suffix || ""}`.trim()}
                  </td>
                  <td>{emp.birthdate}</td>
                  <td>{emp.room}</td>
                  <td>
                    {emp.admissionDate
                      ? new Date(emp.admissionDate).toLocaleString()
                      : ""}
                  </td>
                  <td>
                    <button onClick={() => openEdit(emp)} className="edit-btn">
                      Edit
                    </button>
                    <button
                      onClick={() => viewPatient(emp)}
                      className="view-btn"
                    >
                      View
                    </button>
                    <button
                      onClick={() => deleteEmployee(emp.id)}
                      className="delete-btn"
                    >
                      Delete
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}

export default EmployeeProfiling;
