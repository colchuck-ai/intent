SKILL_SRC := intent
SKILL_DEST := .claude/skills/intent

.PHONY: sync-skill
sync-skill:
	rm -rf $(SKILL_DEST)
	mkdir -p $(SKILL_DEST)
	cp -R $(SKILL_SRC)/ $(SKILL_DEST)/
	@echo "Synced $(SKILL_SRC)/ -> $(SKILL_DEST)/"
