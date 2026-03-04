import React, { useState, useEffect, useCallback } from "react";
import axios from "axios";
import "./PhysicalExamination.css";

/*
  Physical Examination Tabs (UI sections)
  - Skin
  - HEENT
  - Chest        <-- connected to DB table `tsekap_lib_chest` (dynamic)
  - Heart        <-- connected to DB table `tsekap_lib_heart` (dynamic)
  - Abdomen
  - Neurological
  - Digital Rectal
  - Genitourinary

  Dynamic libraries (Chest/Heart) are loaded from backend endpoints:
  - GET /api/lib/chest
  - GET /api/lib/heart

  The dynamic items are normalized and sorted by numeric ID (1..999).
*/

const API_BASE = "http://localhost:8080";

function PhysicalExamination({ patientId }) {
  // State for General Survey
  const [generalSurvey, setGeneralSurvey] = useState("awake");

  // State for Remarks
  const [remarks, setRemarks] = useState("");

  // State for Blood Type
  const [bloodType, setBloodType] = useState("A+");

  // State for active examination category
  const [activeExam, setActiveExam] = useState("Skin");

  // State for all examination findings
  const [examFindings, setExamFindings] = useState({
    // Skin findings
    skin: {
      essentiallyNormal: false,
      weakPulses: false,
      clubbing: false,
      coldClammy: false,
      cyanosis: false,
      edema: false,
      decreasedMobility: false,
      paleNailbeds: false,
      poorSkinTurgor: false,
      rashes: false,
      others: false,
      othersText: "",
    },
    // HEENT findings
    heent: {
      anictericSclerae: false,
      exudates: false,
      essentiallyNormal: false,
      abnormalPupillaryReaction: false,
      cervicalLymphadenopathy: false,
      dryMucousMembrane: false,
      ictericSclerae: false,
      paleConjunctivae: false,
      sunkenEyeballs: false,
      sunkenFontanelle: false,
      intactTympanicMembrane: false,
      pupilsBriskyReactiveToLight: false,
      tonsillopharyngealCongestion: false,
      hypertrophicTonsils: false,
      alarFlaring: false,
      nasalDischarge: false,
      auralDischarge: false,
      palpableMass: false,
      others: false,
      othersText: "",
    },
    // Chest findings
    chest: {
      symmetricalChestExpansion: false,
      lumpsOverBreast: false,
      clearBreathSounds: false,
      retractions: false,
      crackles: false,
      wheezes: false,
      essentiallyNormal: false,
      asymmetricalChestExpansion: false,
      decreasedBreathSounds: false,
      enlargeAxillaryLymphNodes: false,
      others: false,
      othersText: "",
    },
    // Heart findings
    heart: {
      adynamicPrecordium: false,
      normalRateRegularRhythm: false,
      heaves: false,
      murmurs: false,
      essentiallyNormal: false,
      displacedApexBeat: false,
      irregularRhythm: false,
      muffledHeartSounds: false,
      pericardialBulge: false,
      others: false,
      othersText: "",
    },
    // Abdomen findings
    abdomen: {
      flat: false,
      hyperactiveBowelSounds: false,
      palpableMass: false,
      tympanitic: false,
      uterineContraction: false,
      flabby: false,
      globullar: false,
      muscleGuarding: false,
      tenderness: false,
      palpableMass2: false,
      essentiallyNormal: false,
      abdominalRigidity: false,
      abdominalTenderness: false,
      others: false,
      othersText: "",
    },
    // Neurological findings
    neurological: {
      developmentalDelay: false,
      abnormalReflex: false,
      poorAlteredMemory: false,
      poorMuscleTone: false,
      poorCoordination: false,
      seizures: false,
      normal: false,
      motorDeficit: false,
      sensoryDeficit: false,
      essentiallyNormal: false,
      abnormalGait: false,
      abnormalPositionSense: false,
      abnormalSensation: false,
      others: false,
      othersText: "",
    },
    // Digital Rectal findings
    digitalRectal: {
      notApplicable: false,
      essentiallyNormal: false,
      enlargeProstate: false,
      mass: false,
      hemorrhoids: false,
      pus: false,
      others: false,
      othersText: "",
    },
    // Genitourinary findings
    genitourinary: {
      essentiallyNormal: false,
      bloodStainedInExamFinger: false,
      cervicalDilatation: false,
      presenceOfAbnormalDischarge: false,
      others: false,
      othersText: "",
    },
  });

  // === DYNAMIC LIB OPTIONS (fetched from backend) ===
  // `heartOptions` comes from GET /api/lib/heart
  const [heartOptions, setHeartOptions] = useState([]);
  // `chestOptions` comes from GET /api/lib/chest
  const [chestOptions, setChestOptions] = useState([]);
  // `skinOptions` comes from GET /api/lib/skin
  const [skinOptions, setSkinOptions] = useState([]);
  // `heentOptions` comes from GET /api/lib/heent
  const [heentOptions, setHeentOptions] = useState([]);
  // `abdomenOptions` comes from GET /api/lib/abdomen
  const [abdomenOptions, setAbdomenOptions] = useState([]);
  // `neuroOptions` comes from GET /api/lib/neuro
  const [neuroOptions, setNeuroOptions] = useState([]);
  // `digitalRectalOptions` comes from GET /api/lib/digital_rectal
  const [digitalRectalOptions, setDigitalRectalOptions] = useState([]);
  // `genitourinaryOptions` comes from GET /api/lib/genitourinary
  const [genitourinaryOptions, setGenitourinaryOptions] = useState([]);

  useEffect(() => {
    // Generic loader for library endpoints moved inside useEffect to avoid hook-deps warnings
    const loadLibrary = async (url, keyPrefix, setOptions) => {
      try {
        const res = await fetch(url);
        const data = await res.json();

        const normalized = (data || [])
          .map((it) => {
            const rawId = it.id || it.code || it.HEART_ID || it.CHEST_ID || "";
            const id = String(rawId);
            const desc = (it.desc || it.HEART_DESC || it.CHEST_DESC || it.heart_desc || it.chest_desc || "").toString().trim();
            const numId = Number(id.replace(/[^0-9]/g, "")) || 0;
            return { id, desc, numId };
          })
          .sort((a, b) => a.numId - b.numId);

        const opts = normalized.map((it) => {
          const lower = (it.desc || "").toLowerCase();
          if (lower === "other" || lower === "others" || lower.startsWith("other")) {
            return { id: it.id, label: it.desc, key: "others", col: 1 };
          }
          return { id: it.id, label: it.desc, key: `${keyPrefix}${String(it.id)}`, col: 1 };
        });

        // ensure 'others' option (if present) always appears last
        const othersIndex = opts.findIndex((o) => o.key === "others");
        if (othersIndex >= 0) {
          const othersItem = opts.splice(othersIndex, 1)[0];
          opts.push(othersItem);
        }

        setOptions(opts);
      } catch (err) {
        console.warn(`Failed to load ${url}:`, err);
      }
    };

    // use relative URLs so CRA proxy can handle backend routing and avoid CORS
    loadLibrary("/api/lib/heart", "heart_", setHeartOptions);
    loadLibrary("/api/lib/chest", "chest_", setChestOptions);
    loadLibrary("/api/lib/abdomen", "abdomen_", setAbdomenOptions);
    loadLibrary("/api/lib/neuro", "neuro_", setNeuroOptions);
    loadLibrary("/api/lib/digital_rectal", "digitalRectal_", setDigitalRectalOptions);
    loadLibrary("/api/lib/genitourinary", "genitourinary_", setGenitourinaryOptions);
    loadLibrary("/api/lib/skin", "skin_", setSkinOptions);
    loadLibrary("/api/lib/heent", "heent_", setHeentOptions);
  }, []);

  // ========== LOAD / SAVE STATE + ACTIONS ==========
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [saveMsg, setSaveMsg] = useState("");

  const loadData = useCallback(async () => {
    if (!patientId) {
      setLoading(false);
      return;
    }
    setLoading(true);
    try {
      // load general
      const genRes = await axios.get(`${API_BASE}/api/patients/${patientId}/physical-exam/general`);
      const gen = genRes.data || {};
      setGeneralSurvey(gen.general_survey || "awake");
      setRemarks(gen.remarks || "");
      setBloodType(gen.blood_type || "A+");

      // load findings
      const findRes = await axios.get(`${API_BASE}/api/patients/${patientId}/physical-exam/findings`);
      const dbFindings = findRes.data || [];

      // rebuild examFindings structure from current template then apply DB rows
      setExamFindings((prev) => {
        const rebuilt = {};
        Object.keys(prev).forEach((cat) => {
          rebuilt[cat] = {};
          Object.keys(prev[cat]).forEach((field) => {
            rebuilt[cat][field] = field === "othersText" ? "" : false;
          });
        });

        dbFindings.forEach((row) => {
          const cat = row.category;
          const code = row.finding_code;
          if (!rebuilt[cat]) return;
          if (code === "others") {
            rebuilt[cat].others = true;
            rebuilt[cat].othersText = row.others_text || "";
          } else {
            rebuilt[cat][code] = !!row.is_checked;
          }
        });

        return rebuilt;
      });
    } catch (err) {
      console.error("Error loading physical exam:", err);
    } finally {
      setLoading(false);
    }
  }, [patientId]);

  useEffect(() => {
    loadData();
  }, [loadData]);

  const handleSave = async () => {
    if (!patientId) {
      alert("No patient ID. Open this tab from inside a patient record.");
      return;
    }

    setSaving(true);
    setSaveMsg("");
    try {
      // save general
      await axios.post(`${API_BASE}/api/patients/${patientId}/physical-exam/general`, {
        general_survey: generalSurvey,
        remarks: remarks,
        blood_type: bloodType,
      });

      // save findings per category
      for (const catKey of Object.keys(examFindings)) {
        const catFindings = examFindings[catKey] || {};
        // Send all findings (include unchecked ones) so backend can store explicit 0/1
        const findings = Object.keys(catFindings)
          .filter((k) => k !== "othersText")
          .map((k) => {
            if (k === "others") {
              return { finding_code: "others", finding_desc: "Others", is_checked: !!catFindings.others, others_text: catFindings.othersText || "" };
            }
            return { finding_code: k, finding_desc: k, is_checked: !!catFindings[k], others_text: "" };
          });

        await axios.post(`${API_BASE}/api/patients/${patientId}/physical-exam/findings`, { category: catKey, findings });
      }

      setSaveMsg("✓ Saved successfully!");
      setTimeout(() => setSaveMsg(""), 3000);
    } catch (err) {
      const serverMsg = err.response?.data || err.message;
      console.error("PE save failed | patientId:", patientId, "| error:", serverMsg);
      setSaveMsg(`❌ Save failed: ${typeof serverMsg === "string" ? serverMsg : JSON.stringify(serverMsg)}`);
    } finally {
      setSaving(false);
    }
  };

  const handleCheckboxChange = (category, field) => {
    setExamFindings((prev) => ({
      ...prev,
      [category]: {
        ...prev[category],
        [field]: !prev[category][field],
      },
    }));
  };

  const handleOthersTextChange = (category, value) => {
    setExamFindings((prev) => ({
      ...prev,
      [category]: {
        ...prev[category],
        othersText: value,
      },
    }));
  };

  // Examination categories configuration
  const examCategories = [
    { id: "Skin", label: "Skin", key: "skin" },
    { id: "HEENT", label: "HEENT", key: "heent" },
    { id: "Chest", label: "Chest", key: "chest" },
    { id: "Heart", label: "Heart", key: "heart" },
    { id: "Abdomen", label: "Abdomen", key: "abdomen" },
    { id: "Neurological", label: "Neurological", key: "neurological" },
    { id: "Digital Rectal", label: "Digital Rectal", key: "digitalRectal" },
    { id: "Genitourinary", label: "Genitourinary", key: "genitourinary" },
  ];

  // Render examination content based on active category
  const renderExamContent = () => {
    const category = examCategories.find((cat) => cat.id === activeExam);
    if (!category) return null;

    const findings = examFindings[category.key];

    // Define checkbox options for each category
    const checkboxOptions = {
      skin: [
        ...(skinOptions.length > 0
          ? skinOptions
          : [
              { key: "essentiallyNormal", label: "Essentially normal", col: 1 },
              { key: "weakPulses", label: "Weak pulses", col: 1 },
              { key: "clubbing", label: "Clubbing", col: 1 },
              { key: "coldClammy", label: "Cold clammy", col: 1 },
              { key: "cyanosis", label: "Cyanosis/mottled skin", col: 1 },
              { key: "edema", label: "Edema/swelling", col: 1 },
              { key: "decreasedMobility", label: "Decreased mobility", col: 1 },
              { key: "paleNailbeds", label: "Pale nailbeds", col: 1 },
              { key: "poorSkinTurgor", label: "Poor skin turgor", col: 1 },
              { key: "rashes", label: "Rashes/Petechiae", col: 1 },
            ]),
      ],
      heent: [
        ...(heentOptions.length > 0
          ? heentOptions
          : [
              { key: "anictericSclerae", label: "Anicteric sclerae", col: 1 },
              { key: "exudates", label: "Exudates", col: 1 },
              { key: "essentiallyNormal", label: "Essentially Normal", col: 1 },
              {
                key: "abnormalPupillaryReaction",
                label: "Abnormal pupillary reaction",
                col: 1,
              },
              {
                key: "cervicalLymphadenopathy",
                label: "Cervical lympadenopathy",
                col: 1,
              },
              { key: "dryMucousMembrane", label: "Dry mucous membrane", col: 1 },
              { key: "ictericSclerae", label: "Icteric sclerae", col: 1 },
              { key: "paleConjunctivae", label: "Pale conjunctivae", col: 1 },
              { key: "sunkenEyeballs", label: "Sunken eyeballs", col: 1 },
              { key: "sunkenFontanelle", label: "Sunken fontanelle", col: 1 },
              {
                key: "intactTympanicMembrane",
                label: "Intact tympanic mebrane",
                col: 1,
              },
              {
                key: "pupilsBriskyReactiveToLight",
                label: "Pupils briskly reactive to light",
                col: 1,
              },
              {
                key: "tonsillopharyngealCongestion",
                label: "Tonsillopharyngeal congestion",
                col: 2,
              },
              { key: "hypertrophicTonsils", label: "Hypertrophic tonsils", col: 2 },
              { key: "alarFlaring", label: "Alar flaring", col: 2 },
              { key: "nasalDischarge", label: "Nasal discharge", col: 2 },
              { key: "auralDischarge", label: "Aural discharge", col: 2 },
              { key: "palpableMass", label: "Palpable mass", col: 2 },
            ]),
      ],
      // Chest checkbox options — prefers dynamic `chestOptions` from backend when available
      chest: [
        ...(chestOptions.length > 0
          ? chestOptions
          : [
              {
                key: "symmetricalChestExpansion",
                label: "Symmetrical chest expansion",
                col: 1,
              },
              { key: "lumpsOverBreast", label: "Lumps over breast(s)", col: 1 },
              { key: "clearBreathSounds", label: "Clear breath sounds", col: 1 },
              { key: "retractions", label: "Retractions", col: 1 },
              { key: "crackles", label: "Crackles/rales", col: 1 },
              { key: "wheezes", label: "Wheezes", col: 1 },
              { key: "essentiallyNormal", label: "Essentially normal", col: 1 },
              {
                key: "asymmetricalChestExpansion",
                label: "Asymmetrical chest expansion",
                col: 1,
              },
              {
                key: "decreasedBreathSounds",
                label: "Decreased breath sounds",
                col: 1,
              },
              {
                key: "enlargeAxillaryLymphNodes",
                label: "Enlarge Axillary Lymph Nodes",
                col: 1,
              },
            ]),
      ],
      // Heart checkbox options — prefers dynamic `heartOptions` from backend when available
      heart: [
        ...(heartOptions.length > 0
          ? heartOptions
          : [
              { key: "adynamicPrecordium", label: "Adynamic precordium", col: 1 },
              {
                key: "normalRateRegularRhythm",
                label: "Normal rate regular rhythm",
                col: 1,
              },
              { key: "heaves", label: "Heaves/trills", col: 1 },
              { key: "murmurs", label: "Murmurs", col: 1 },
              { key: "essentiallyNormal", label: "Essentially normal", col: 1 },
              { key: "displacedApexBeat", label: "Displaced apex beat", col: 1 },
              { key: "irregularRhythm", label: "Irregular rhythm", col: 1 },
              { key: "muffledHeartSounds", label: "Muffled heart sounds", col: 1 },
              { key: "pericardialBulge", label: "Pericardial bulge", col: 1 },
            ]),
      ],
      abdomen: [
        ...(abdomenOptions.length > 0
          ? abdomenOptions
          : [
              { key: "flat", label: "Flat", col: 1 },
              {
                key: "hyperactiveBowelSounds",
                label: "Hyperactive bowel sounds",
                col: 1,
              },
              { key: "palpableMass", label: "Palpable mass(es)", col: 1 },
              { key: "tympanitic", label: "Tympanitic/dull abdomen", col: 1 },
              { key: "uterineContraction", label: "Uterine contraction", col: 1 },
              { key: "flabby", label: "Flabby", col: 1 },
              { key: "globullar", label: "Globullar", col: 1 },
              { key: "muscleGuarding", label: "Muscle guarding", col: 1 },
              { key: "tenderness", label: "Tenderness", col: 1 },
              { key: "palpableMass2", label: "Palpable mass", col: 1 },
              { key: "essentiallyNormal", label: "Essentially normal", col: 1 },
              { key: "abdominalRigidity", label: "Abdominal rigidity", col: 1 },
              { key: "abdominalTenderness", label: "Abdominal tenderness", col: 2 },
            ]),
      ],
      neurological: [
        ...(neuroOptions.length > 0
          ? neuroOptions
          : [
              { key: "developmentalDelay", label: "Developmental delay", col: 1 },
              { key: "abnormalReflex", label: "Abnormal reflex(es)", col: 1 },
              { key: "poorAlteredMemory", label: "Poor/altered memory", col: 1 },
              { key: "poorMuscleTone", label: "Poor muscle tone/strength", col: 1 },
              { key: "poorCoordination", label: "Poor coordination", col: 1 },
              { key: "seizures", label: "Seizures", col: 1 },
              { key: "normal", label: "Normal", col: 1 },
              { key: "motorDeficit", label: "Motor Deficit", col: 1 },
              { key: "sensoryDeficit", label: "Sensory Deficit", col: 1 },
              { key: "essentiallyNormal", label: "Essentially normal", col: 1 },
              { key: "abnormalGait", label: "Abnormal gait", col: 1 },
              {
                key: "abnormalPositionSense",
                label: "Abnormal position sense",
                col: 1,
              },
              { key: "abnormalSensation", label: "Abnormal sensation", col: 2 },
            ]),
      ],
      digitalRectal: [
        ...(digitalRectalOptions.length > 0
          ? digitalRectalOptions
          : [
              { key: "notApplicable", label: "Not Applicable", col: 1 },
              { key: "essentiallyNormal", label: "Essentially normal", col: 1 },
              { key: "enlargeProstate", label: "Enlarge Prostate", col: 1 },
              { key: "mass", label: "Mass", col: 1 },
              { key: "hemorrhoids", label: "Hemorrhoids", col: 1 },
              { key: "pus", label: "Pus", col: 1 },
            ]),
      ],
      genitourinary: [
        ...(genitourinaryOptions.length > 0
          ? genitourinaryOptions
          : [
              { key: "essentiallyNormal", label: "Essentially normal", col: 1 },
              {
                key: "bloodStainedInExamFinger",
                label: "Blood stained in exam finger",
                col: 1,
              },
              { key: "cervicalDilatation", label: "Cervical dilatation", col: 1 },
              {
                key: "presenceOfAbnormalDischarge",
                label: "Presence of abnormal discharge",
                col: 1,
              },
            ]),
      ],
    };

    const options = checkboxOptions[category.key] || [];
    const col1Options = options.filter((opt) => opt.col === 1);
    const col2Options = options.filter((opt) => opt.col === 2);

    const hasDbOther = (category.key === "heart" && heartOptions.some((o) => o.key === "others")) ||
      (category.key === "chest" && chestOptions.some((o) => o.key === "others")) ||
      (category.key === "skin" && skinOptions.some((o) => o.key === "others")) ||
      (category.key === "heent" && heentOptions.some((o) => o.key === "others")) ||
      (category.key === "abdomen" && abdomenOptions.some((o) => o.key === "others")) ||
      (category.key === "neurological" && neuroOptions.some((o) => o.key === "others")) ||
      (category.key === "digitalRectal" && digitalRectalOptions.some((o) => o.key === "others")) ||
      (category.key === "genitourinary" && genitourinaryOptions.some((o) => o.key === "others"));

    return (
      <div className="pe-content-details">
        <div className="pe-header-row">
          <h4>
            {category.id === "HEENT"
              ? "HEENT (Head, Eyes, Ears, Nose, Throat)"
              : category.id === "Skin"
              ? "Skin / Extremities"
              : category.label}
          </h4>
        </div>

        <div className="pe-grid-layout">
          {col1Options.length > 0 && (
            <div className="pe-checkboxes">
              {col1Options.map((option) => {
                const inputId = `${category.key}-${option.key}`;
                return (
                  <label key={option.key} className="checkbox-item" htmlFor={inputId}>
                    <input
                      id={inputId}
                      type="checkbox"
                      checked={findings[option.key] || false}
                      onChange={() => handleCheckboxChange(category.key, option.key)}
                    />
                    <span>{option.label}</span>
                  </label>
                );
              })}
            </div>
          )}

          {col2Options.length > 0 && (
            <div className="pe-checkboxes">
              {col2Options.map((option) => {
                const inputId = `${category.key}-${option.key}`;
                return (
                  <label key={option.key} className="checkbox-item" htmlFor={inputId}>
                    <input
                      id={inputId}
                      type="checkbox"
                      checked={findings[option.key] || false}
                      onChange={() => handleCheckboxChange(category.key, option.key)}
                    />
                    <span>{option.label}</span>
                  </label>
                );
              })}
            </div>
          )}
        </div>

        <div className="pe-others-area">
          {!hasDbOther && (
            <label className="checkbox-item" htmlFor={`${category.key}-others`}>
              <input
                id={`${category.key}-others`}
                type="checkbox"
                checked={findings.others || false}
                onChange={() => handleCheckboxChange(category.key, "others")}
              />
              <span>Others</span>
            </label>
          )}

          <textarea
            value={findings.othersText || ""}
            onChange={(e) => handleOthersTextChange(category.key, e.target.value)}
            placeholder="lorem ipsum"
            disabled={!findings.others}
          />
        </div>

        {/* Save Button + Status */}
        <div className="pe-section">
          <button className="pe-save-btn" onClick={handleSave} disabled={saving}>
            {saving ? "Saving..." : "💾 Save"}
          </button>
          {saveMsg && (
            <div className={`pe-save-msg ${saveMsg.startsWith("✓") ? "success" : "error"}`}>
              {saveMsg}
            </div>
          )}
        </div>

      </div>

        
    );
  };

  if (loading) return (
    <div className="loading-container">Loading physical examination data...</div>
  );

  if (!patientId) return (
    <div className="error-container">
      ⚠️ No patient ID. Make sure PatientView passes <code>patientId={"{id}"}</code> to &lt;PhysicalExamination&gt;.
    </div>
  );

  return (
    <div className="physical-exam-layout">
      {/* LEFT PANE */}
      <div className="pe-left-pane">
        {/* General Survey Section */}
        <div className="pe-section">
          <h4>1. General Survey</h4>
          <div className="radio-group-vertical">
            <label>
              <input
                type="radio"
                name="generalSurvey"
                value="awake"
                checked={generalSurvey === "awake"}
                onChange={(e) => setGeneralSurvey(e.target.value)}
              />
              Awake and Alert
            </label>
            <label>
              <input
                type="radio"
                name="generalSurvey"
                value="altered"
                checked={generalSurvey === "altered"}
                onChange={(e) => setGeneralSurvey(e.target.value)}
              />
              Altered Sensorium
            </label>
          </div>
        </div>

        {/* Remarks Section */}
        <div className="pe-section">
          <h4>2. Remarks</h4>
          <textarea
            className="remarks-box"
            value={remarks}
            onChange={(e) => setRemarks(e.target.value)}
            placeholder="lorem ipsum 12121"
          />
        </div>

        {/* Blood Type Section */}
        <div className="pe-section bloodtype-section">
          <h4>Bloodtype</h4>
          <div className="bloodtype-badge">1. Patient Blood Type</div>
          <div className="radio-group-vertical compact">
            {['A+', 'B+', 'AB+', 'O+', 'A-', 'B-', 'AB-', 'O-'].map((type) => (
              <label key={type}>
                <input
                  type="radio"
                  name="bloodType"
                  value={type}
                  checked={bloodType === type}
                  onChange={(e) => setBloodType(e.target.value)}
                />
                {type}
              </label>
            ))}
          </div>
        </div>
      </div>

      {/* RIGHT PANE */}
      <div className="pe-right-pane">
        {/* Navigation Tabs */}
        <div className="pe-nav">
          {examCategories.map((cat) => (
            <button
              key={cat.id}
              className={`pe-nav-link ${activeExam === cat.id ? "active" : ""}`}
              onClick={() => setActiveExam(cat.id)}
            >
              {cat.label}
            </button>
          ))}
        </div>

        {/* Dynamic Content */}
        {renderExamContent()}
      </div>
    </div>
  );
}

export default PhysicalExamination;
