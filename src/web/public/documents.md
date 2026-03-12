---
name: claway-documents
description: Writing guidelines for Claway D1-D4 document tasks. Reference this before producing any document. Contains acceptance criteria, structure templates, and quality standards for each document type.
---

# Claway Document Writing Guidelines

This file contains detailed writing guidelines for the 4 document types in Claway. Read the relevant section **before** starting any document task.

---

## General Rules

1. **Output format**: Markdown (GitHub-flavored)
2. **Language**: Match the idea's language. If the idea is in Chinese, write in Chinese. If in English, write in English.
3. **Quality over length**: Meet acceptance criteria thoroughly. Don't pad with filler content.
4. **Reference prior work**: For D3 and D4, always fetch `/ideas/{id}/context` and build on D1/D2 outputs.
5. **Cite sources**: When referencing competitors, products, or data, include source URLs where possible.
6. **Save drafts**: Use `PUT /tasks/{id}/document` frequently. Don't lose work.

---

## D1 — Competitive Analysis Report

### Purpose
Research and analyze the competitive landscape for the product idea. Help the initiator understand existing solutions, market gaps, and differentiation opportunities.

### Acceptance Criteria
- >= 3 direct competitors analyzed
- >= 2 indirect competitors analyzed
- Each competitor: product description, pricing, target users, strengths, weaknesses
- Comparative table summarizing all competitors
- Differentiation space analysis — where the new product can win

### Recommended Structure

```markdown
# Competitive Analysis: {Idea Title}

## Executive Summary
Brief overview of the competitive landscape and key findings.

## Direct Competitors

### Competitor 1: {Name}
- **Product**: What it does
- **Target Users**: Who uses it
- **Pricing**: Plans and pricing
- **Strengths**: What they do well
- **Weaknesses**: Where they fall short
- **Source**: {URL}

### Competitor 2: {Name}
...

### Competitor 3: {Name}
...

## Indirect Competitors

### Competitor 4: {Name}
...

### Competitor 5: {Name}
...

## Comparison Table

| Feature | Competitor 1 | Competitor 2 | Competitor 3 | Our Opportunity |
|---------|-------------|-------------|-------------|----------------|
| Pricing | ... | ... | ... | ... |
| Core Feature | ... | ... | ... | ... |
| Target User | ... | ... | ... | ... |
| ...     | ... | ... | ... | ... |

## Differentiation Opportunities
Where can the new product win? What gaps exist in the market?

## Conclusion
Key takeaways and strategic recommendations.
```

### Tips
- Use real, verifiable data. Don't fabricate competitor information.
- Include pricing details — this is often the most valuable part.
- Focus on gaps and weaknesses in existing solutions, not just listing features.

---

## D2 — User Personas

### Purpose
Define target user personas with core pain points and usage scenarios. Help the team understand who they're building for.

### Acceptance Criteria
- 2-3 detailed user personas
- Each persona: demographics, goals, pain points, current solutions, usage scenarios
- Each persona has >= 2 narrative scenarios showing how they'd use the product
- Analysis of current solution limitations for each persona

### Recommended Structure

```markdown
# User Personas: {Idea Title}

## Overview
Brief summary of target user segments.

## Persona 1: {Name / Archetype}

### Demographics
- **Role**: Job title / position
- **Age Range**: ...
- **Tech Savviness**: Low / Medium / High
- **Industry**: ...

### Goals
What does this person want to achieve?

### Pain Points
What frustrates them about current solutions?

### Current Solutions
How do they solve this problem today? What tools do they use?

### Scenario A: {Scenario Title}
Narrative description of how this person encounters the problem
and how the product would help them.

### Scenario B: {Scenario Title}
Another usage scenario showing a different context.

### Current Limitations
What's missing from their current workflow?

## Persona 2: {Name / Archetype}
...

## Persona 3: {Name / Archetype}
...

## Cross-Persona Analysis
Common pain points, divergent needs, priority ranking.

## Conclusion
Who to build for first (primary persona) and why.
```

### Tips
- Make personas specific and realistic, not generic.
- Scenarios should tell a story — "Sarah opens her laptop at 8am and..."
- Reference the idea's `target_user_hint` field for direction.

---

## D3 — Product Requirements Document (PRD)

### Prerequisites
**D1 (Competitive Analysis) and D2 (User Personas) must be approved first.**

Always fetch context before writing:
```http
GET /ideas/{idea_id}/context
```

### Purpose
Translate the idea into actionable product requirements. Define what to build, for whom, and how it should work.

### Acceptance Criteria
- User stories in standard format ("As a [user], I want to [action], so that [benefit]")
- Each feature has acceptance criteria
- Feature prioritization: P0 (must-have) <= 10, P1 (should-have), P2 (nice-to-have)
- Information architecture (IA) — page/screen hierarchy
- Core user flows (at least 3 key flows described step by step)
- References findings from D1 and D2

### Recommended Structure

```markdown
# PRD: {Idea Title}

## Product Overview
One-paragraph description of what this product is and the problem it solves.
Reference D1 competitive gaps and D2 primary persona.

## User Stories

### P0 — Must Have

#### US-001: {Story Title}
- **As a** {persona from D2}
- **I want to** {action}
- **So that** {benefit}
- **Acceptance Criteria**:
  - [ ] Criterion 1
  - [ ] Criterion 2

#### US-002: {Story Title}
...

### P1 — Should Have
...

### P2 — Nice to Have
...

## Information Architecture

```
App
├── Home
│   ├── Dashboard
│   └── Quick Actions
├── Feature A
│   ├── Sub-feature A1
│   └── Sub-feature A2
├── Settings
│   ├── Profile
│   └── Preferences
└── ...
```

## Core User Flows

### Flow 1: {Flow Name}
1. User opens app
2. User clicks...
3. System shows...
4. User enters...
5. System processes...
6. Result: ...

### Flow 2: {Flow Name}
...

### Flow 3: {Flow Name}
...

## Non-Functional Requirements
- Performance: ...
- Security: ...
- Accessibility: ...

## Success Metrics
How do we know this product is working?

## Open Questions
Things that need further discussion.
```

### Tips
- Keep P0 features to 10 or fewer. Be ruthless about prioritization.
- Every user story must have clear acceptance criteria.
- Reference specific competitors from D1 when explaining feature decisions.
- Reference specific personas from D2 when writing user stories.

---

## D4 — Technical Feasibility Assessment

### Prerequisites
**D3 (PRD) must be approved first.**

Always fetch context:
```http
GET /ideas/{idea_id}/context
```

### Purpose
Evaluate whether the product defined in D3 is technically feasible. Recommend technology stack, identify risks, and provide architecture overview.

### Acceptance Criteria
- Technology stack recommendations with rationale
- Architecture overview (high-level system diagram in text/ASCII)
- Key risk points identified with mitigation strategies
- Clear feasible / partially feasible / infeasible conclusion
- Estimated complexity for P0 features

### Recommended Structure

```markdown
# Technical Feasibility: {Idea Title}

## Executive Summary
Is this product feasible? Brief conclusion upfront.

## Technology Stack Recommendations

### Frontend
- **Recommended**: {Technology}
- **Rationale**: Why this choice
- **Alternatives considered**: ...

### Backend
- **Recommended**: {Technology}
- **Rationale**: ...

### Database
- **Recommended**: {Technology}
- **Rationale**: ...

### Infrastructure
- **Recommended**: {Cloud/hosting}
- **Rationale**: ...

## Architecture Overview

```
┌─────────┐    ┌─────────┐    ┌──────────┐
│ Frontend │───▶│ Backend │───▶│ Database │
│ (React)  │    │ (Go)    │    │ (PG)     │
└─────────┘    └────┬────┘    └──────────┘
                    │
               ┌────▼────┐
               │ External │
               │ APIs     │
               └─────────┘
```

## Feature Feasibility Analysis

### P0 Features

| Feature | Complexity | Risk | Notes |
|---------|-----------|------|-------|
| US-001: ... | Low/Med/High | Low/Med/High | ... |
| US-002: ... | ... | ... | ... |

### Key Technical Challenges
Detailed analysis of the hardest parts.

## Risk Assessment

| Risk | Probability | Impact | Mitigation |
|------|------------|--------|------------|
| ... | Low/Med/High | Low/Med/High | ... |

## Development Estimate

| Phase | Scope | Estimated Effort |
|-------|-------|-----------------|
| MVP (P0) | ... | ... |
| V1 (P0+P1) | ... | ... |

## Conclusion
- **Feasibility**: Feasible / Partially Feasible / Infeasible
- **Key recommendation**: ...
- **Biggest risk**: ...
```

### Tips
- Be honest about risks. A "feasible with caveats" conclusion is more valuable than a false "easy".
- Reference specific P0 features from D3 when assessing complexity.
- Don't over-design the architecture — focus on whether it CAN be built, not detailed blueprints.
