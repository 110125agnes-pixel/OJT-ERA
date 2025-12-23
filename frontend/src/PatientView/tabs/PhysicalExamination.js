import React, { useState } from "react";
import "./PhysicalExamination.css";

function PhysicalExamination() {
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
      ],
      heent: [
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
          label: "Pupils brisky reactive to light",
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
      ],
      chest: [
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
      ],
      heart: [
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
      ],
      abdomen: [
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
      ],
      neurological: [
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
      ],
      digitalRectal: [
        { key: "notApplicable", label: "Not Applicable", col: 1 },
        { key: "essentiallyNormal", label: "Essentially normal", col: 1 },
        { key: "enlargeProstate", label: "Enlarge Prostate", col: 1 },
        { key: "mass", label: "Mass", col: 1 },
        { key: "hemorrhoids", label: "Hemorrhoids", col: 1 },
        { key: "pus", label: "Pus", col: 1 },
      ],
      genitourinary: [
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
      ],
    };

    const options = checkboxOptions[category.key] || [];
    const col1Options = options.filter((opt) => opt.col === 1);
    const col2Options = options.filter((opt) => opt.col === 2);

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
              {col1Options.map((option) => (
                <label key={option.key} className="checkbox-item">
                  <input
                    type="checkbox"
                    checked={findings[option.key] || false}
                    onChange={() =>
                      handleCheckboxChange(category.key, option.key)
                    }
                  />
                  <span>{option.label}</span>
                </label>
              ))}
            </div>
          )}

          {col2Options.length > 0 && (
            <div className="pe-checkboxes">
              {col2Options.map((option) => (
                <label key={option.key} className="checkbox-item">
                  <input
                    type="checkbox"
                    checked={findings[option.key] || false}
                    onChange={() =>
                      handleCheckboxChange(category.key, option.key)
                    }
                  />
                  <span>{option.label}</span>
                </label>
              ))}
            </div>
          )}
        </div>

        <div className="pe-others-area">
          <label className="checkbox-item">
            <input
              type="checkbox"
              checked={findings.others || false}
              onChange={() => handleCheckboxChange(category.key, "others")}
            />
            <span>Others</span>
          </label>
          <textarea
            value={findings.othersText || ""}
            onChange={(e) =>
              handleOthersTextChange(category.key, e.target.value)
            }
            placeholder="lorem ipsum"
            disabled={!findings.others}
          />
        </div>
      </div>
    );
  };

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
            {["A+", "B+", "AB+", "O+", "A-", "B-", "AB-", "O-"].map((type) => (
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
