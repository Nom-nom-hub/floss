# Documentation Format Standards

This document defines the standards for all FLOSS documentation to ensure consistency and quality across all documents.

## Document Structure

All documentation files should follow this structure:

1. **Document Title** - H1 header with the document title
2. **Brief Description** - 1-2 sentence description of what the document covers
3. **Table of Contents** - For longer documents (>100 lines)
4. **Content Sections** - Organized with appropriate H2 and H3 headers
5. **Related Links** - Links to related documentation at the end

## Writing Style

### Voice and Tone
- Use second person ("you") rather than first person or third person
- Use active voice rather than passive voice
- Be clear and concise
- Avoid jargon unless necessary, and explain when used

### Language Guidelines
- Use American English
- Use present tense
- Use consistent terminology throughout
- Define acronyms on first use

### Formatting Standards

#### Headers
- Use H1 for document titles
- Use H2 for main sections
- Use H3 for subsections
- Do not skip header levels

#### Code Examples
- Use triple backticks with language identifier for code blocks
- Use single backticks for inline code
- Include comments in complex examples
- Use realistic examples that users can copy and modify

#### Lists
- Use bulleted lists for unordered items
- Use numbered lists for ordered steps
- Keep list items parallel in structure
- Use consistent punctuation in list items

#### Links
- Use descriptive link text rather than "click here"
- Link to related documentation where appropriate
- Use relative links for internal documentation
- Verify all links are working

## Metadata

Each document should include metadata at the top:

```markdown
---
title: Document Title
description: Brief description of the document
last_updated: YYYY-MM-DD
audience: [user|developer|maintainer]
---
```

## Diagrams

For diagrams, prefer:
1. ASCII art for simple diagrams
2. Mermaid.js for complex diagrams
3. Link to external images for screenshots

## File Naming

- Use lowercase filenames
- Use hyphens to separate words
- Use descriptive names
- Use .md extension for all documentation files

## Versioning

Documentation should be kept up to date with the latest release. When features are added or changed:

1. Update documentation in the same pull request when possible
2. For major changes, create a separate documentation update PR
3. Mark deprecated features as such
4. Include version information for new features