#include "context.h"
#include "malloc.h"
#include "../implementations/constant.h"

context *newContext(void* goContext) {
    context *output = (context*)malloc(sizeof(context));
    output->goContext = goContext;
    output->remainingMagnitudeData = (uint16_t*)malloc((MAGNITUDE_LONG_MSG_SIZE + PREAMBULE_BIT_SIZE) * sizeof(uint16_t));
    output->remainingMagnitudeLengthByte = 0;

    // fprintf(stderr, "Allocating memory: %ld\n", (MAGNITUDE_LONG_MSG_SIZE + PREAMBULE_BIT_SIZE) * sizeof(uint16_t));
    return output;
}