#ifndef __RTL_SDR_H__
#define __RTL_SDR_H__

#include <stdlib.h>
#include <unistd.h>
#include "rtl-sdr.h"


extern uint16_t magnitude[129*129];


extern int goRtlsrdData(unsigned char *buf, uint32_t len, void *ctx);

int rtlsdrReadAsync(rtlsdr_dev_t *dev, void *ctx, uint32_t buf_num, uint32_t buf_len);
void rtlsdrProcessRaw(unsigned char *buf, uint32_t len, void *ctx);
void initTables(int debug);
int decodeMessage(uint16_t* magnitudeBuffer, char * message);

#endif