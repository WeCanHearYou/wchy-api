# How to localize Fider

English files `en/client.json` and `en/server.json` are the only files that should be modified during development.

All other locale files are sourced from [Crowdin](https://crowdin.com/project/fider). If you're looking to contribute with a new translation, a correction or improvement, please use Crowdin.

# How to sync Git with Crowdin

*Note:* This section is a how-to process for Fider admins only.
*Note 2:* Would be great to automate this...

## Source Strings

1. Visit https://crowdin.com/project/fider/settings#files
2. For each file, click Update and upload the respective file

## Sync English from Git to Crowdin

1. Visit https://crowdin.com/project/fider/en#
2. For each file, select Upload translations
   - Mark `Allow target translation to match source`
   - Mark `Approve added translations`

## Sync other locales from Crowdin to Git

1. Review and approve translations
2. Visit https://crowdin.com/project/fider/settings#translations
3. Select build and download
4. Copy files to locale folder 
5. Update percentages on `locales.ts` based on Crowdin numbers
6. Create PR