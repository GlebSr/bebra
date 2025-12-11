TEMPLATE_FILE=template.config.yaml
OUTPUT_FILE=config.yaml

generate:
	@if [ ! -f "$(TEMPLATE_FILE)" ]; then \
  		echo "Template file not found: $(TEMPLATE_FILE)"; \
  		exit 1; \
	fi
	envsubst < $(TEMPLATE_FILE) > $(OUTPUT_FILE)
	@echo "Generated $(OUTPUT_FILE) from $(TEMPLATE_FILE)"

clean:
	@if [ -f "$(OUTPUT_FILE)" ]; then \
  		rm -f $(OUTPUT_FILE); \
  		echo "Removed $(OUTPUT_FILE)"; \
	else \
  		echo "$(OUTPUT_FILE) not found, nothing to remove."; \
	fi
